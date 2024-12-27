package main

import (
	"digitalcorporation/pkg/utils"
	"net/http"
)

// static resources
func (app *application) css(w http.ResponseWriter, r *http.Request) {
	bytes, err := utils.ReadFile(r.URL.Path)
	if err != nil {
		app.serverError(w, err)
		return
	}

	mw := app.minifyWriter("text/css", w)
	defer mw.Close()

	app.writeStaticResource(w, r, bytes, "text/css; charset=utf-8", mw)
}

func (app *application) js(w http.ResponseWriter, r *http.Request) {
	bytes, err := utils.ReadFile(r.URL.Path)
	if err != nil {
		app.serverError(w, err)
		return
	}

	mw := app.minifyWriter("application/javascript", w)
	defer mw.Close()

	app.writeStaticResource(w, r, bytes, "application/javascript; charset=utf-8", mw)
}

func (app *application) img(w http.ResponseWriter, r *http.Request) {
	bytes, err := utils.ReadFile(r.URL.Path)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// get mimeType of requested resource (detection for image work good)
	mimeType := http.DetectContentType(bytes)
	app.writeStaticResource(w, r, bytes, mimeType, nil)
}

// endpoints
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	res, err := app.services.imageUploader.Upload(r, "files")
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, res)
}

func (app *application) auth(w http.ResponseWriter, r *http.Request) {
	if userID, ok := r.Context().Value(app.constants.contextUserID).(string); ok {
		app.infoLog.Println("userID: ", userID)
	}

	app.writeJSON(w, http.StatusOK, map[string]string{
		"message": "authenticated",
	})
}
