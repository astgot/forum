package main

import (
	"log"

	"github.com/astgot/forum/internal/server"
)

func main() {
	config := server.NewConfig() // generating config for server
	server := server.New(config) // creating new instance based on config
	// Starting the Server
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
