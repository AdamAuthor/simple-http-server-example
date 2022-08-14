package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"kbtu_go_6/internal/http"
	"kbtu_go_6/internal/store/postgres"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	url := os.Getenv("URL")

	store := postgres.NewDB()
	if err := store.Connect(url); err != nil {
		panic(err)
	}
	defer store.Close()

	//Creating and run new server
	srv := http.NewServer(context.Background(), ":8080", store)
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}

	srv.WaitForGT()
}
