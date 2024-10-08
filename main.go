package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Fenroe/chirpy/internal/config"
	"github.com/Fenroe/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	fmt.Println(os.Getenv("PLATFORM"))
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("An error occurred: %v", err)
	}
	dbQueries := database.New(db)
	apiConfig := config.Config{}
	apiConfig.Queries = dbQueries
	handler := http.NewServeMux()
	fileServer := apiConfig.MiddlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	handler.HandleFunc("GET /api/healthz", apiConfig.Readiness)
	handler.HandleFunc("GET /admin/metrics", apiConfig.GetMetrics)
	handler.HandleFunc("POST /admin/reset", apiConfig.ResetMetrics)
	handler.HandleFunc("POST /api/users", apiConfig.CreateUser)
	handler.HandleFunc("POST /api/chirps", apiConfig.CreateChirp)
	handler.HandleFunc("GET /api/chirps", apiConfig.GetChirps)
	handler.Handle("/app/", fileServer)
	server := http.Server{
		Handler: handler,
		Addr:    ":8080",
	}
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
