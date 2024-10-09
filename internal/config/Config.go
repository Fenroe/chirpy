package config

import (
	"sync/atomic"

	"github.com/Fenroe/chirpy/internal/database"
)

type Config struct {
	FileserverHits atomic.Int32
	Queries        *database.Queries
	JWTSecret      string
	PulkaKey       string
}
