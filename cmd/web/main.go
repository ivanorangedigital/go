package main

import (
	"database/sql"
	"digitalcorporation/pkg/models/mysql"
	"digitalcorporation/pkg/services"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
)

type _constants struct {
	// env
	envJwtSecret    string
	envCookieSecret string

	// cookie
	cookieJwt string

	// context
	contextUserID string
}

type _services struct {
	imageUploader *services.ImageUploader
}

type _models struct {
	userModel *mysql.UserModel
}

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	constants     *_constants
	services      *_services
	models        *_models
	templateCache map[string]*template.Template
	minifyWriter  func(mediatype string, w io.Writer) io.WriteCloser
}

func main() {
	// flags
	addr := flag.String("addr", "localhost:3000", "Set address port")
	dsn := flag.String("dsn", "digital:Orange521!@/ecommerce?parseTime=true", "Data source name")
	flag.Parse()

	// create custom loggers
	infoLog := log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR ", log.Ldate|log.Ltime|log.Lshortfile)

	// load env file
	err := loadEnv()
	if err != nil {
		errorLog.Fatal(err)
	}
	// ---

	// instance db connection
	db, err := openDB(*dsn)

	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()
	// ---

	// set constants for applictaion
	constants := &_constants{
		// env
		envJwtSecret:    os.Getenv("JWT_SECRET"),
		envCookieSecret: os.Getenv("COOKIE_SECRET"),

		// cookie
		cookieJwt: "JWT_TOKEN",

		// context
		contextUserID: "USER_ID",
	}

	// set services for application
	services := &_services{
		imageUploader: &services.ImageUploader{
			RootDir:           "./ui/static/img/",
			PrefixRequest:     "/img/",
			MaxLength:         2,
			MaxSizeFile:       2 << 20,
			AllowedExtensions: []string{"png"},
		},
	}

	// set models for application
	models := &_models{
		userModel: &mysql.UserModel{DB: db},
	}

	// set cached templates
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// set minify for application (html, css, js)
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("application/javascript", js.Minify)

	// initialize applications struct
	app := &application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		constants:     constants,
		services:      services,
		models:        models,
		templateCache: templateCache,
		minifyWriter:  m.Writer,
	}

	// create server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Server starting on port %s", strings.Split(*addr, ":")[1])
	errorLog.Fatal(srv.ListenAndServe())
}

func loadEnv() error {
	bytes, err := os.ReadFile(".env")
	if err != nil {
		return err
	}

	data := string(bytes)
	list := strings.Split(data, "\n")

	for _, item := range list {
		// skip empty line
		if item == "" {
			continue
		}

		splittedItem := strings.Split(item, "=")
		if len(splittedItem) != 2 {
			return fmt.Errorf("Invalid .env file")
		}

		key := splittedItem[0]
		value := splittedItem[1]

		// remove " or '
		firstChar := value[0]
		lastChar := value[len(value)-1]
		if (firstChar == '"' && lastChar == '"') || (firstChar == '\'' && lastChar == '\'') {
			value = value[1 : len(value)-1]
		}

		os.Setenv(key, value)
	}

	return nil
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
