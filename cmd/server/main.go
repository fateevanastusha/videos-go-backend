package main

import (
	"log"
	"os"

	"github.com/fateevanastusha/videos-go-backend/internal/server"
)

func main() {
	address := "localhost:3005"
	if port := os.Getenv("PORT"); port != "" {
		address = ":" + port
	}
	if err := server.Start(address); err != nil {
		log.Fatal(err)
	}
}
