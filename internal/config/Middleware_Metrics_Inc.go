package config

import "net/http"

func (C *Config) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		C.FileserverHits.Add(1)
		next.ServeHTTP(res, req)
	})
}
