package main

import (
	"digitalcorporation/pkg/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
)

func (app *application) chainMiddleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	// reverse cycle for better user usage, the first handler passed will be first execute
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func (app *application) writeHTML(w http.ResponseWriter, data any, name string) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("Page with %s name doesn't exists", name))
		return
	}

	mw := app.minifyWriter("text/html", w)
	defer mw.Close()

	// write header
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	// execute template set
	if err := ts.Execute(mw, data); err != nil {
		app.serverError(w, err)
	}
}

func (app *application) writeStaticResource(w http.ResponseWriter, r *http.Request, bytes []byte, mimeType string, mw io.Writer) {
	lastModified, err := utils.LastModifiedDateFile(r.URL.Path)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// write header
	w.Header().Set("Last-Modified", lastModified.Format(http.TimeFormat))
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", utils.ConvertDaysInSeconds(1)))
	w.Header().Set("Content-Type", mimeType)
	w.WriteHeader(http.StatusOK)

	// write data
	if mw != nil {
		mw.Write(bytes)
	} else {
		w.Write(bytes)
	}
}

func (app *application) writeCookie(w http.ResponseWriter, key, value string, maxAge uint64) error {
	cookies := w.Header()["Set-Cookie"]

	for _, cookie := range cookies {
		// use +1 to avoid problems with cookie key like ciao and setted cookie key is ciaosono
		// in this case cookieKey will be ciaos, that different from ciao=
		// otherwise the value would be the same ciao (error)
		cookieKey := cookie[:len(key)+1]

		if cookieKey == key+"=" {
			return fmt.Errorf("Cookie key '%s' already setted", cookieKey)
		}
	}

	// encrypt cookie value
	value, err := utils.Encrypt(value, app.constants.envCookieSecret)
	if err != nil {
		return err
	}

	w.Header().Add("Set-Cookie", fmt.Sprintf("%s=%s; HttpOnly; SameSite=Strict; Secure; Max-Age=%d; Path=/;", key, value, maxAge))
	return nil
}

func (app *application) getCookie(r *http.Request, key string) (string, error) {
	cookies := r.Cookies()

	for _, cookie := range cookies {
		if cookie.Name == key {

			// decrypt cookie value
			return utils.Decrypt(cookie.Value, app.constants.envCookieSecret)
		}
	}

	return "", nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, v any) {
	bytes, err := json.Marshal(v)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// set header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(bytes)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

// v must be a pointer
func (app *application) readJSON(r *http.Request, v any) error {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(v); err != nil {
		return err
	}

	return nil
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(trace)
	app.writeJSON(w, http.StatusInternalServerError, map[string]string{
		"message": http.StatusText(http.StatusInternalServerError),
	})
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	app.writeJSON(w, status, map[string]string{
		"message": http.StatusText(status),
	})
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) unauthorized(w http.ResponseWriter) {
	app.clientError(w, http.StatusUnauthorized)
}
