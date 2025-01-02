package about

import (
	"digitalcorporation/cmd/web2/handlers"
	"net/http"
)

func init() {
	handlers.NewHandler().RegisterRoute("/about", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the About page!"))
	})
}
