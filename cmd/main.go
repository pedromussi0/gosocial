package main

import (
	"log"
	"net/http"
)

func main() {
	api := &api{addr: ":8080"}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /users", api.getUsersHandler)
	mux.HandleFunc("GET /users/:id", api.getUsersHandler)
	mux.HandleFunc("POST /users", api.createUsersHandler)
	// mux.HandleFunc("PUT /users/:id", api.updateUsersHandler)
	// mux.HandleFunc("DELETE /users/:id", api.deleteUsersHandler)

	srv := &http.Server{
		Addr:    api.addr,
		Handler: mux,
	}

	log.Printf("Server starting on %s", api.addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
