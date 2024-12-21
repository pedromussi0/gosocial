package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{
		addr: os.Getenv("ADDR"),
	}

	app := &application{
		config: cfg,
	}

	// Mount routes and start server
	app.mount()
	log.Fatal(app.run())
}
