package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/pedromussi0/gosocial.git/internal/db"
	"github.com/pedromussi0/gosocial.git/internal/store"
)

const version = "0.0.1"

//	@title	GoSocial API

//	@description	API for GoSocial, a social network.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath						/v1
//
// @securityDefinitions.api_key	ApiKeyAuth
// @in								header
// @name							Authorization
// @description
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{
		addr:   os.Getenv("ADDR"),
		apiURL: os.Getenv("API_URL"),
		db: dbConfig{
			addr:         os.Getenv("DB_ADDR"),
			maxOpenConns: 30,
			maxIdleConns: 30,
			maxIdleTime:  "15m",
		},
		env: os.Getenv("ENV"),
	}

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Println("database connection pool established")

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	// Mount routes and start server
	app.mount()
	log.Fatal(app.run())
}
