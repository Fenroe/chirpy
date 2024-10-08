package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Fenroe/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits  atomic.Int32
	databaseQueries *database.Queries
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(res, req)
	})
}

func (cfg *apiConfig) getNumbrOfRequests(res http.ResponseWriter, _ *http.Request) {
	res.Header().Add("Content-Type", "text/html")
	res.WriteHeader(200)
	htmlBody := fmt.Sprintf(`
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>`, cfg.fileserverHits.Load())
	res.Write([]byte(htmlBody))
}

func (cfg *apiConfig) resetNumberOfRequests(res http.ResponseWriter, _ *http.Request) {
	defer res.Header().Add("Content-Type", "text/html")
	mode := os.Getenv("PLATFORM")
	if mode != "dev" {
		res.WriteHeader(403)
		htmlBody := `
		<html>
			<body>
				<h1>Forbidden/h1>
				<p>You don't have permission to do that</p>
			</body>
		</html>`
		res.Write([]byte(htmlBody))
	} else {
		cfg.databaseQueries.DeleteUsers(context.Background())
		res.WriteHeader(200)
		previousValue := cfg.fileserverHits.Swap(0)
		htmlBody := fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>You've successfully reset the visit count from %d to %d</p>
			</body>
		</html>`, previousValue, cfg.fileserverHits.Load())
		res.Write([]byte(htmlBody))
	}
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	fmt.Println(os.Getenv("PLATFORM"))
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("An error occurred: %v", err)
	}
	dbQueries := database.New(db)
	config := apiConfig{}
	config.databaseQueries = dbQueries
	handler := http.NewServeMux()
	fileServer := config.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	handler.HandleFunc("GET /api/healthz", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(200)
		res.Write([]byte("OK"))
	})
	handler.HandleFunc("GET /admin/metrics", config.getNumbrOfRequests)
	handler.HandleFunc("POST /admin/reset", config.resetNumberOfRequests)
	handler.HandleFunc("POST /api/validate_chirp", validateChirp)
	handler.HandleFunc("POST /api/users", config.createUser)
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
