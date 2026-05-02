package main

import (
	"fmt"
	"log"
	"main/db"
	"main/routes"
	"net/http"
	"os"
)

func main() {
	database, err := db.Init()
	if err != nil {
		log.Fatal("[main] Error initializing database: ", err)
	}
	defer database.Close()

	if os.Getenv("SEED_DB") == "true" {
		if err := db.Seed(database); err != nil {
			log.Fatal("[main] Error seeding database: ", err)
		}
	}

	r := routes.SetupRoutes(database)

	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
