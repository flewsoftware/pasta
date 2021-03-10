package pasta_core

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"github.com/kardianos/osext"
	_ "github.com/mattn/go-sqlite3"
	"pasta/filecrypt-core"
	"pasta/pasta_core/db"
)

const DBFileName = "pass.pasta_db"

func GetDBLocation() (string, error) {
	s, err := osext.ExecutableFolder()
	if err != nil {
		return "", err
	}
	return s + "/" + DBFileName, nil
}

func GenerateDB(masterPass string) error {
	location, err := GetDBLocation()
	if err != nil {
		return err
	}

	database, err := sql.Open("sqlite3", location)
	if err != nil {
		return err
	}

	// create secrets table
	statement, err := database.Prepare(db.CreateSecretsTableIfNotExist)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}

	// create hash store table
	statement, err = database.Prepare(db.CreateHashStoreTableIfNotExist)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}

	// add master password hash
	statement, err = database.Prepare(db.AddMasterPasswordHash)
	if err != nil {
		return err
	}
	hash, err := HashMasterPassword(masterPass)
	if err != nil {
		return err
	}
	_, err = statement.Exec(hash)
	if err != nil {
		return fmt.Errorf("%s (The database might be already initalized)", err.Error())
	}

	return nil
}

func AddSecretToDB(secretName string, secret string, masterPassword string) error {
	if err, val := IsMasterPasswordCorrect(masterPassword); err != nil {
		return fmt.Errorf("%s (Database might not be initalized)", err.Error())
	} else if val == false {
		return fmt.Errorf("pasta: Invalid master password")
	}

	location, err := GetDBLocation()
	if err != nil {
		return err
	}

	database, _ := sql.Open("sqlite3", location)

	statement, _ := database.Prepare(db.AddSecret)
	secretE, err := EncryptSecret(secret, masterPassword)
	if err != nil {
		return err
	}

	_, err = statement.Exec(secretName, secretE)
	return err
}

func GetSecretFromDB(ExceptName string, masterPassword string) (string, error) {
	err, secrets := GetSecretsFromDB(masterPassword)
	if err != nil {
		return "", err
	}

	for secretName, val := range secrets {
		if secretName == ExceptName {
			return val, nil
		}
	}

	return "", nil
}

func GetSecretsFromDB(masterPassword string) (error, map[string]string) {
	location, err := GetDBLocation()
	if err != nil {
		return err, nil
	}

	database, _ := sql.Open("sqlite3", location)

	err, res := IsMasterPasswordCorrect(masterPassword)
	if err != nil {
		return err, nil
	} else if !res {
		return fmt.Errorf("pasta: invalid master password"), nil
	}

	passwords := make(map[string]string)
	rows, _ := database.Query(db.QueryAllSecrets)
	var keyName string
	var encryptedValue []byte
	for rows.Next() {
		rows.Scan(&keyName, &encryptedValue)
		p, err := DecryptSecret(encryptedValue, masterPassword)
		if err != nil {
			return err, nil
		}
		passwords[keyName] = string(p)
	}

	return nil, passwords
}

func EncryptSecret(password string, masterPassword string) (filecrypt.EncryptedData, error) {
	data, err := filecrypt.EncryptSHA256([]byte(password), filecrypt.Passphrase(masterPassword))
	return data, err
}

func DecryptSecret(encryptedPassword []byte, masterPassword string) (filecrypt.DecryptedData, error) {
	data, err := filecrypt.DecryptSHA256(encryptedPassword, filecrypt.Passphrase(masterPassword))
	return data, err
}

func GetMasterPasswordHash() ([]byte, error) {
	location, err := GetDBLocation()
	if err != nil {
		return nil, err
	}

	database, err := sql.Open("sqlite3", location)
	if err != nil {
		return nil, err
	}
	rows, err := database.Query(db.QueryMasterPasswordHash)
	if err != nil {
		return nil, err
	}
	var (
		keyName string
		blob    []byte
	)
	for rows.Next() {
		rows.Scan(&keyName, &blob)
		if keyName == db.MasterPasswordHashKey {
			return blob, nil
		}
	}
	return nil, nil
}

func IsMasterPasswordCorrect(masterPasswordToCheck string) (error, bool) {
	p, err := GetMasterPasswordHash()
	if err != nil {
		return err, false
	}

	checkHash, err := HashMasterPassword(masterPasswordToCheck)
	if err != nil {
		return err, false
	}
	if string(checkHash) == string(p) {
		return nil, true
	}
	return nil, false
}

func HashMasterPassword(p string) (filecrypt.Hash, error) {
	var salt []byte
	// Generate a Salt
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	return filecrypt.CreateHashArgon(p, salt)
}
