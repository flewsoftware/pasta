package db

const (
	MasterPasswordHashKey = "mp"
)

const (
	CreateSecretsTableIfNotExist   = "CREATE TABLE IF NOT EXISTS secrets (keyName TEXT PRIMARY KEY, encryptedValue BLOB)"
	CreateHashStoreTableIfNotExist = "CREATE TABLE IF NOT EXISTS hashStore (keyName TEXT PRIMARY KEY, hash BLOB)"
)

const (
	AddMasterPasswordHash = "INSERT INTO hashStore (keyName, hash) VALUES ('" + MasterPasswordHashKey + "', ?)"
	AddSecret             = "INSERT INTO secrets (keyName, encryptedValue) VALUES (?, ?)"
)

const (
	QueryAllSecrets         = "SELECT keyName, encryptedValue FROM secrets"
	QueryMasterPasswordHash = "SELECT keyName, hash FROM hashStore"
)
