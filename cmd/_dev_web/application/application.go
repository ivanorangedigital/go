package application

import (
	"sync"
)

var (
	instance *application
	once     sync.Once
)

type application struct {
}

func NewApplication() *application {
	once.Do(func() {
		instance = &application{}
	})
	return instance
}
