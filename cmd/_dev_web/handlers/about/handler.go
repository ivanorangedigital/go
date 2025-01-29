package about

import (
	"digitalcorporation/cmd/web/handlers"
	"net/http"
)

func init() {
	handlers.NewHandlers().RegisterRoute("/about", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the About page!"))
	})
}
