package home

import (
	"digitalcorporation/cmd/web/handlers"
	"net/http"
)

func init() {
	handlers.NewHandlers().RegisterRoute("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Home page!"))
	})
}
