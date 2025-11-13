package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	gonanoid "github.com/matoous/go-nanoid/v2"

	sqlc "github.com/abdulyazidi/urlshortner/generated"
	_ "modernc.org/sqlite"
)

func main() {
	fmt.Println("Hello, World!")
	ctx := context.Background()

	db, err := sql.Open("sqlite", "./db/urlshortner.db")
	if err != nil {
		log.Fatal("Error connecting to db", err)
	}
	queries := sqlc.New(db)

	app := &App{queries: queries}

	url := "https://google.com"
	shortenedURL, err := gonanoid.New(8)

	if err != nil {
		log.Fatal("Couldn't generate a nanoid", err)
	}

	// insert shortened url
	result, err := queries.CreateURL(ctx, sqlc.CreateURLParams{ID: shortenedURL, OriginalUrl: url})

	if err != nil {
		log.Fatal("Failed to insert new shortened url", err)
	} else {
		fmt.Println("CreateURL result: ", result)
	}

	// retrieve  original url
	originalURL, err := queries.GetOriginalURL(ctx, result.ID)

	if err != nil {
		log.Println("Failed to query original url", err)
	}

	fmt.Printf("Shortened url: %v\nOriginal url: %v\n", shortenedURL, originalURL)

	http.HandleFunc("/{$}", HomepageHandler)
	http.HandleFunc("/", app.RedirectToOriginalURLHandler)
	http.HandleFunc("/shorten", app.ShortenURLHandler)
	http.ListenAndServe(":8090", nil)
	fmt.Println("Exiting...")
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

func (app *App) ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusInternalServerError)
	}
	shortURL, err := gonanoid.New(8)
	if err != nil {
		fmt.Println("Error")
	}
	fmt.Println(shortURL)
	fmt.Fprintf(w, "Method: %s\n", r.Method)
}

func HomepageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	fmt.Fprintf(w, "Homepage\n")
}

type App struct {
	queries *sqlc.Queries
}
