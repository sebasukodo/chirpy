package handler

import (
	"sync/atomic"

	"github.com/sebasukodo/chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	DbQueries      *database.Queries
	Platform       string
	TokenSecret    string
	PolkaApiKey    string
}
