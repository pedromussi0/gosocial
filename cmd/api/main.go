package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/pedromussi0/gosocial.git/internal/store"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{
		addr: os.Getenv("ADDR"),
	}

	store := store.NewStorage(nil)

	app := &application{
		config: cfg,
		store:  store,
	}

	// Mount routes and start server
	app.mount()
	log.Fatal(app.run())
}
