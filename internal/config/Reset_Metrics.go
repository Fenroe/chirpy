package config

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

func (C *Config) ResetMetrics(res http.ResponseWriter, _ *http.Request) {
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
		C.Queries.DeleteUsers(context.Background())
		res.WriteHeader(200)
		previousValue := C.FileserverHits.Swap(0)
		htmlBody := fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>You've successfully reset the visit count from %d to %d</p>
			</body>
		</html>`, previousValue, C.FileserverHits.Load())
		res.Write([]byte(htmlBody))
	}
}
