package main

import (
	"log"

	_ "github.com/joho/godotenv/autoload" // ¯\_(ツ)_/¯
)

func main() {
	config := newConfig()
	log.Printf("Listening on port: %d", config.Port)
	initHTTPServer(config)
}