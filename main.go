package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/jpleatherland/blogaggregator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	PORT := os.Getenv("PORT")
	DB_CONN_STRING := os.Getenv("DB_CONN_STRING")
	db, err := sql.Open("postgres", DB_CONN_STRING)
	if err != nil {
		log.Fatalf("unable to open db connection: %v", err)
	}
	dbQueries := database.New(db)
	resources := Resources{
		DB: dbQueries,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthz", healthCheck)
	mux.HandleFunc("GET /v1/err", errorCheck)

	mux.HandleFunc("POST /v1/users", resources.createUser)
	mux.HandleFunc("GET /v1/users", resources.authMiddleware(resources.getUser))

	mux.HandleFunc("POST /v1/feeds", resources.authMiddleware(resources.createFeed))
	mux.HandleFunc("GET /v1/feeds", resources.getAllFeeds)

	mux.HandleFunc("POST /v1/feed_follows", resources.authMiddleware(resources.createFeedFollow))
	mux.HandleFunc("GET /v1/feed_follows", resources.authMiddleware(resources.getFeedFollowsByUserId))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", resources.authMiddleware(resources.deleteFeedFollow))

	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: mux,
	}

	log.Printf("Listening on port: " + PORT)
	log.Fatal(server.ListenAndServe())
}

func healthCheck(rw http.ResponseWriter, _ *http.Request) {
	type payload struct {
		Status string `json:"status"`
	}
	response := payload{
		Status: "ok",
	}
	respondWithJSON(rw, 200, response)
}

func errorCheck(rw http.ResponseWriter, _ *http.Request) {
	respondWithError(rw, 500, "internal server error")
}

type Resources struct {
	DB *database.Queries
}
