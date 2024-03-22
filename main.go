package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/nljackson2020/boot-dev-blog-aggregator/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	dbURL := os.Getenv("CONN")
	if dbURL == "" {
		log.Fatal("$CONN must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	config := &apiConfig{
		DB: dbQueries,
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/readiness", handlerReadiness)
	v1Router.Get("/err", handlerError)

	v1Router.Get("/users", config.middlewareAuth(config.handlerUserGet))
	v1Router.Post("/users", config.handlerUserCreate)

	v1Router.Get("/feeds", config.handlerGetFeeds)
	v1Router.Post("/feeds", config.middlewareAuth(config.handlerFeedCreate))

	v1Router.Get("/feed_follows", config.middlewareAuth(config.handlerFeedFollowsGet))
	v1Router.Post("/feed_follows", config.middlewareAuth(config.handlerFeedFollowCreate))
	v1Router.Delete("/feed_follows/{feedFollowID}", config.middlewareAuth(config.handlerFeedFollowDelete))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	const collectionConcurrency = 10
	const collectionInterval = time.Minute
	go startScraping(dbQueries, collectionConcurrency, collectionInterval)

	log.Printf("Serving files on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())

}
