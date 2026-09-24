package main

import (
	"log"

	"github.com/fateevanastusha/videos-go-backend/internal/apiserver"
)

func main() {
	address := "localhost:3005"
	if err := apiserver.Start(address); err != nil {
		log.Fatal(err)
	}
}
