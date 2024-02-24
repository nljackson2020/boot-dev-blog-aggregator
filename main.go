package main

import (
	//"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	//	"github.com/nljackson2020/boot-dev-blog-aggregator/internal/database"

	_ "github.com/lib/pq"
)

// type apiConfig struct {
// 	DB *database.Queries
// }

func main() {
	const filepathRoot = "."
	godotenv.Load()

	port := os.Getenv("PORT")
	//dbURL := os.Getenv("DB_URL")

	//db, err := sql.Open("postgres", dbURL)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//dbQueries := database.New(db)

	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	r.Use(cors.Handler)

	v1Router := chi.NewRouter()
	v1Router.Get("/readiness", handlerReadiness)
	v1Router.Get("/err", handlerError)
	v1Router.Post("/users", handlerCreateUser)

	r.Mount("/v1", v1Router)

	corsMux := middlewareCors(r)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

}
