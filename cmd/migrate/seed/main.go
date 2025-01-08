package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pedromussi0/gosocial.git/internal/db"
	"github.com/pedromussi0/gosocial.git/internal/store"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := os.Getenv("DB_ADDR")

	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store, conn)
}
