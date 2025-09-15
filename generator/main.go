package main

import (
	"fmt"
	"log"

	"app/database"

	config "bitbucket.org/sadeemTechnology/backend-config"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	var tableName string
	fmt.Print("Enter the table name: ")
	fmt.Scanln(&tableName)

	cfg := config.GetSettings()
	db, err := database.OpenSQLX(cfg.DSN)
	if err != nil {
		log.Fatalf("couldn't open db: %s", err.Error())
		return
	}
	log.Println("database connection pool established")

	if tableName == "" {
		fmt.Println("table name cannot be empty.")
		return
	}

	if err := generateModelFile(tableName, db); err != nil {
		fmt.Printf("Error generating controller file: %s\n", err)
		return
	}
	if err := generateControllerFile(tableName); err != nil {
		fmt.Printf("Error generating controller file: %s\n", err)
		return
	}
}
