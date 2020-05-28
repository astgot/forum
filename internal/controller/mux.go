package controller

import (
	"net/http"

	"github.com/astgot/forum/internal/database"
)

// Multiplexer ....
type Multiplexer struct {
	Mux *http.ServeMux
	db  *database.Database
}

// NewMux ...
func NewMux() *Multiplexer {
	return &Multiplexer{
		Mux: http.NewServeMux(),
		db:  database.NewDB(database.NewConfig()),
	}
}
