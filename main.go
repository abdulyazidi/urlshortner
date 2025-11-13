package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	gonanoid "github.com/matoous/go-nanoid/v2"

	sqlc "github.com/abdulyazidi/urlshortner/generated"
	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "./db/urlshortner.db")
	if err != nil {
		log.Fatal("Error connecting to db", err)
	}
	queries := sqlc.New(db)
	app := &App{queries: queries}
	app.InitTempl()
	// Register handlers
	http.HandleFunc("/{$}", app.HomepageHandler)
	http.HandleFunc("/", app.RedirectToOriginalURLHandler)
	var port int = 8090
	fmt.Printf("Running on port: %d\nhttp://localhost:%d", port, port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	fmt.Println("Exiting...")
}

type App struct {
	queries *sqlc.Queries
	tmpl    *template.Template
}

type PageData struct {
	ShortenedURL string
	Error        string
}

func (app *App) InitTempl() {
	app.tmpl = template.Must(template.ParseFiles("./html/homepage.html"))
}

func (app *App) RedirectToOriginalURLHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	shortURL := r.URL.Path[1:]

	originalURL, err := app.queries.GetOriginalURL(ctx, shortURL)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	fmt.Println(r.URL.Path)
	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

func (app *App) HomepageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var data PageData

	// Handle POST (form submission)
	if r.Method == http.MethodPost {
		ctx := r.Context()

		if err := r.ParseForm(); err != nil {
			data.Error = "Failed to parse form: " + err.Error()
		} else {
			originalURL := r.FormValue("url")
			shortURL, err := gonanoid.New(8)
			if err != nil {
				data.Error = "Failed to generate short URL: " + err.Error()
			} else {
				app.queries.CreateURL(ctx, sqlc.CreateURLParams{ID: shortURL, OriginalUrl: originalURL})
				data.ShortenedURL = shortURL
			}
		}
	}

	// Render template (for both GET and POST)
	if err := app.tmpl.Execute(w, data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Printf("Template execution error: %v\n", err)
	}
}
