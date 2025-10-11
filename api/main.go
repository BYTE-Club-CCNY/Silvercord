package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	grpcClient "github.com/BYTE-Club-CCNY/Silvercord/api/grpc"
	"github.com/BYTE-Club-CCNY/Silvercord/api/routes"
	"github.com/joho/godotenv"
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

	// Initialize gRPC client for LLM service
	llmClient, err3 := grpcClient.NewLLMClient("localhost:50051")
	if err3 != nil {
		log.Fatal("Error creating LLM gRPC client: ", err3)
	}
	defer llmClient.Close()

	r := routes.SetupRoutes(client, llmClient)

	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
