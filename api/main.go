package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"main/routes"
	"net/http"
	"os"

	"github.com/supabase-community/supabase-go"
)

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	supabaseUrl := os.Getenv("SUPABASE_URL_PRIVATE")
	supabaseAnonKey := os.Getenv("SUPABASE_KEY_PRIVATE")
	client, err2 := supabase.NewClient(
		supabaseUrl,
		supabaseAnonKey,
		&supabase.ClientOptions{},
	)
	if err2 != nil {
		log.Fatal("Error creating Supabase client: ", err2)
	}
	r := routes.SetupRoutes(client)

	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
