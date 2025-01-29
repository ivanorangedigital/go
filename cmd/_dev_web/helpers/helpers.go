package helpers

import (
	"net/http"
	"sync"
)

var (
	instance *helpers
	once     sync.Once
)

type helpers struct{}

func (h *helpers) ChainMiddleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	// reverse cycle for better user usage, the first handler passed will be first execute
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func NewHelpers() *helpers {
	once.Do(func() {
		instance = new(helpers)
	})
	return instance
}
