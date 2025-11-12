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

	// Generate model files
	if err := generateModelFiles(tableName, db); err != nil {
		fmt.Printf("Error generating model files: %s\n", err)
		return
	}

	// Generate controller files
	if err := generateControllerFiles(tableName); err != nil {
		fmt.Printf("Error generating controller files: %s\n", err)
		return
	}

	modelName := toPascalCase(tableName)

	if err := updateAPIRoutesFile(modelName); err != nil {
		log.Printf("Warning: Could not update API routes: %v", err)
	} else {
		log.Printf("Updated API routes with %s", modelName)
	}

	if err := updateControllersFile(modelName); err != nil {
		log.Printf("Warning: Could not update controllers: %v", err)
	} else {
		log.Printf("Updated controllers with %s", modelName)
	}

	if err := updateModelsFile(modelName); err != nil {
		log.Printf("Warning: Could not update models: %v", err)
	} else {
		log.Printf("Updated models with %s", modelName)
	}

	log.Printf("Successfully generated model and controller for table: %s", tableName)
}
