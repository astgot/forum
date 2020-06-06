package controller

import (
	"html/template"
	"net/http"

	"github.com/astgot/forum/internal/database"
)

var tpl = template.Must(template.ParseGlob("web/templates/*"))

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
