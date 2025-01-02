package home

import (
	"digitalcorporation/cmd/web2/handlers"
	"net/http"
)

func init() {
	handlers.NewHandler().RegisterRoute("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Home page!"))
	})
}
