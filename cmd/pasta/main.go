package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"io"
	"log"
	"os"
	"pasta/pasta_core"
	"strconv"
	"time"
)

func main() {

	prompt := promptui.Select{
		Label: "Select Mode",
		Items: []string{"get", "add", "init", "reset", "backup"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	switch result {
	case "get":
		var (
			masterPass string
			secretName string
		)

		println("Enter Master Password")
		_, _ = fmt.Scanln(&masterPass)

		println("Enter Secret Name")
		_, _ = fmt.Scanln(&secretName)

		secret, err := pasta_core.GetSecretFromDB(secretName, masterPass)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%s: %s\n", secretName, secret)
		break
	case "add":
		var masterPass string
		var secretName string
		var secret string

		println("Enter Master Password")
		_, _ = fmt.Scanln(&masterPass)

		println("Enter Secret Name")
		_, _ = fmt.Scanln(&secretName)

		println("Enter Secret")
		_, _ = fmt.Scanln(&secret)

		err := pasta_core.AddSecretToDB(secretName, secret, masterPass)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Done!")
	case "init":
		var masterPass string

		println("Enter Master Password")
		_, err = fmt.Scanln(&masterPass)
		if err != nil {
			log.Fatalln(err)
		}

		if err := pasta_core.GenerateDB(masterPass); err != nil {
			log.Fatalln(err)
		}
	case "reset":
		{
			prompt = promptui.Select{
				Label: "Are you sure that you want to delete the pasta password storage database",
				Items: []string{"yes", "no"},
			}
			_, result, err := prompt.Run()
			if err != nil {
				log.Fatalf("Prompt failed %v\n", err)
			}

			if result == "yes" {
				location, err := pasta_core.GetDBLocation()
				if err != nil {
					log.Fatalln(err)
				}

				if err := os.Remove(location); err != nil {
					log.Fatalln(err)
				}
				log.Println("Successfully removed pasta password database use mode 'init' to initialize the database again")
			}
		}
	case "backup":
		{
			location, err := pasta_core.GetDBLocation()
			if err != nil {
				log.Fatalln(err)
			}
			f, err := os.Open(location)
			if err != nil {
				log.Fatalln(err)
			}

			ff, err := os.Create("./passBackup(" + strconv.FormatInt(time.Now().Unix(), 16) + ").pasta_db")
			if err != nil {
				log.Fatalln(err)
			}

			_, err = io.Copy(ff, f)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

}
