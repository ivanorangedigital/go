package handlers

import (
	"digitalcorporation/cmd/web/helpers"
	"errors"
	"net/http"
	"sync"
)

var (
	instance *handlers
	once     sync.Once
)

type handlers struct {
	mux               *http.ServeMux
	globalMiddlewares []func(http.Handler) http.Handler
	mutex             sync.Mutex
}

func (h *handlers) RegisterRoute(path string, handlerFunc http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	h.mux.Handle(path, helpers.NewHelpers().ChainMiddleware(handlerFunc, middlewares...))
}

func (h *handlers) SetGlobalMiddlewares(globalMiddlewares ...func(http.Handler) http.Handler) error {
	// protection for go concurrency
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// check if already setted
	if h.globalMiddlewares != nil {
		return errors.New("This function was already called")
	}

	// set global middlewares
	h.globalMiddlewares = globalMiddlewares
	return nil
}

func (h *handlers) GetMainHandler() http.Handler {
	return helpers.NewHelpers().ChainMiddleware(h.mux, h.globalMiddlewares...)
}

// singletone instance
func NewHandlers() *handlers {
	once.Do(func() {
		instance = &handlers{
			mux: http.NewServeMux(),
		}
	})

	return instance
}
