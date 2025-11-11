package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

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

	fmt.Println("Exiting...")
}
