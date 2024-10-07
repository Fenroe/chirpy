package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(res, req)
	})
}

func (cfg *apiConfig) getNumbrOfRequests(res http.ResponseWriter, _ *http.Request) {
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(200)
	res.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) resetNumberOfRequests(res http.ResponseWriter, _ *http.Request) {
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(200)
	cfg.fileserverHits.Store(0)
	res.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}

func main() {
	config := apiConfig{}
	handler := http.NewServeMux()
	fileServer := config.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	handler.HandleFunc("/healthz", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(200)
		res.Write([]byte("OK"))
	})
	handler.HandleFunc("/metrics", config.getNumbrOfRequests)
	handler.HandleFunc("/reset", config.resetNumberOfRequests)
	handler.Handle("/app/", fileServer)
	server := http.Server{
		Handler: handler,
		Addr:    ":8080",
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
