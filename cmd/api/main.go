package main

import (
	"log"
)

func main() {
	cfg := config{
		addr: ":8080",
	}

	app := &application{
		config: cfg,
	}

	// Mount routes and start server
	app.mount()
	log.Fatal(app.run())
}
