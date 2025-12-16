package main

import (
	"fmt"
	"log"
	"main/routes"
	"net/http"
	"os"

	"github.com/supabase-community/supabase-go"
)

func main() {
	supabaseUrl := os.Getenv("SUPABASE_URL_PRIVATE")
	supabaseAnonKey := os.Getenv("SUPABASE_KEY_PRIVATE")

	if supabaseUrl == "" || supabaseAnonKey == "" {
		log.Fatal("[main] SUPABASE_URL_PRIVATE or SUPABASE_KEY_PRIVATE environment variables not set")
	}

	client, err2 := supabase.NewClient(
		supabaseUrl,
		supabaseAnonKey,
		&supabase.ClientOptions{},
	)
	if err2 != nil {
		log.Fatal("[main] Error creating Supabase client: ", err2)
	}

	r := routes.SetupRoutes(client)

	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
