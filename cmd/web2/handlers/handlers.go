package handlers

import (
	"digitalcorporation/cmd/web2/helpers"
	"errors"
	"net/http"
	"sync"
)

var (
	instance *handler
	once     sync.Once
)

type handler struct {
	mux               *http.ServeMux
	globalMiddlewares []func(http.Handler) http.Handler
	mutex             sync.Mutex
}

func (h *handler) RegisterRoute(path string, handlerFunc http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	h.mux.Handle(path, helpers.ChainMiddleware(handlerFunc, middlewares...))
}

func (h *handler) SetGlobalMiddlewares(globalMiddlewares ...func(http.Handler) http.Handler) error {
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

func (h *handler) GetHandler() http.Handler {
	return helpers.ChainMiddleware(h.mux, h.globalMiddlewares...)
}

// singletone instance
func NewHandler() *handler {
	once.Do(func() {
		instance = &handler{
			mux: http.NewServeMux(),
		}
	})

	return instance
}
