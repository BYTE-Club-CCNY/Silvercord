package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/supabase-community/supabase-go"
)

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	supabaseUrl := os.Getenv("SUPABASE_URL_PRIVATE")
	supabaseAnonKey := os.Getenv("SUPABASE_KEY_PRIVATE")
	_, err = supabase.NewClient(
		supabaseUrl,
		supabaseAnonKey,
		&supabase.ClientOptions{},
	)
	if err != nil {
		log.Fatal("Error creating Supabase client: ", err)
	}
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Chi Server Running"))
	})
	fmt.Println("Listening on port 8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
