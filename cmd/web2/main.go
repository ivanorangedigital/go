package main

import (
	"digitalcorporation/cmd/web2/handlers"
	"log"
	"net/http"

	// -- handlers
	_ "digitalcorporation/cmd/web2/handlers/about"
	_ "digitalcorporation/cmd/web2/handlers/home"
	_ "digitalcorporation/cmd/web2/handlers/newh"
	// -- end handlers
)

func main() {
	// create server
	srv := &http.Server{
		Addr: ":3000",
		// ErrorLog: errorLog,
		Handler: handlers.NewHandler().GetHandler(),
	}

	log.Fatal(srv.ListenAndServe())
}
