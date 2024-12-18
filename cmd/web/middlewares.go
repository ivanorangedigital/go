package main

import (
	"context"
	"digitalcorporation/pkg/utils"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"
)

func (app *application) staticResourcePathValidatorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// clean path and dot checker
		cleanedPath := path.Clean(r.URL.Path[1:])
		if !strings.Contains(cleanedPath, ".") {
			app.notFound(w)
			return
		}

		// concatenate it to base root
		r.URL.Path = fmt.Sprintf("./ui/static/%s", cleanedPath)

		next.ServeHTTP(w, r)
	})
}

func (app *application) staticResourceIsModifiedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lastModified, err := utils.LastModifiedDateFile(r.URL.Path)
		if err != nil {
			app.serverError(w, err)
			return
		}

		// retrieve last date of modified file (in utc)
		if ifModifiedSince := r.Header.Get("If-Modified-Since"); ifModifiedSince != "" {
			modifiedSince, err := http.ParseTime(ifModifiedSince)
			if err != nil {
				app.serverError(w, err)
				return
			}

			// return status not modified if modifiedSince is before lastModified
			if !lastModified.After(modifiedSince) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// request start
		start := time.Now()
		app.infoLog.Printf("Request %s incoming to %s", r.Method, r.URL.Path)

		// next
		next.ServeHTTP(w, r)

		// calc total time
		duration := time.Since(start)
		app.infoLog.Printf("Execution time for %s to %s: %v\n\n", r.Method, r.URL.Path, duration)
	})
}

func (app *application) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Println("Authentication start")

		token, err := app.getCookie(r, app.constants.cookieJwt)

		if err != nil || token == "" {
			app.infoLog.Println("Authentication failed")
			app.unauthorized(w)
			return
		}

		// validate jwt
		userId, err := utils.ValidateJWT(token, app.constants.envJwtSecret)
		if err != nil {
			app.infoLog.Println("Authentication failed")
			app.unauthorized(w)
			return
		}

		// add userId to context for next handlers
		ctx := context.WithValue(r.Context(), app.constants.contextUserID, userId)

		app.infoLog.Printf("Authentication end, passed key %s with value %s in context", app.constants.contextUserID, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
