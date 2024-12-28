package home

import (
	"digitalcorporation/cmd/web2/handlers"
	"net/http"
)

func init() {
	handler := handlers.NewHandler()
	handler.RegisterRoute("/", homeHandler)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Home page!"))
}
