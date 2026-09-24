package main

import (
	"log"

	"github.com/fateevanastusha/videos-go-backend/internal/server"
)

func main() {
	address := "localhost:3005"
	if err := server.Start(address); err != nil {
		log.Fatal(err)
	}
}
