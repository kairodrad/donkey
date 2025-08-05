package main

import (
	"log"

	"github.com/example/donkey/internal/server"
)

// @title Donkey API
// @version 1.0
// @description API for the Donkey card game
// @BasePath /
func main() {
	r := server.New()
	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}
