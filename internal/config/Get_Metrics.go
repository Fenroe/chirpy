package config

import (
	"fmt"
	"net/http"
)

func (C *Config) GetMetrics(res http.ResponseWriter, _ *http.Request) {
	res.Header().Add("Content-Type", "text/html")
	res.WriteHeader(200)
	htmlBody := fmt.Sprintf(`
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>`, C.FileserverHits.Load())
	res.Write([]byte(htmlBody))
}
