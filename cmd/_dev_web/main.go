package main

import (
	// "digitalcorporation/cmd/web/application"
	"digitalcorporation/cmd/web/handlers"
	"log"
	"net/http"

	// -- handlers
	_ "digitalcorporation/cmd/web/handlers/about"
	_ "digitalcorporation/cmd/web/handlers/home"
	// -- end handlers
)

func main() {
	// instance application
	// application := application.NewApplication()

	// instace handlers
	handlers := handlers.NewHandlers()

	// create server
	srv := &http.Server{
		Addr: ":3000",
		// ErrorLog: errorLog,
		Handler: handlers.GetMainHandler(),
	}

	log.Fatal(srv.ListenAndServe())
}
