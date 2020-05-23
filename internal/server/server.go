package server

import (
	"html/template"
	"net/http"

	"github.com/astgot/forum/internal/database"
)

var tpl = template.Must(template.ParseGlob("web/templates/*"))

// Server ..
type Server struct {
	config   *Config
	mux      *http.ServeMux
	database *database.Database
}

// New - generates instance to support service
func New(config *Config) *Server {
	return &Server{
		config: config,
		mux:    http.NewServeMux(),
	}
}

// Start - Initializing server
func (s *Server) Start() error {

	s.ConfigureRouter()
	if err := s.ConfigureDB(); err != nil {
		return err
	}
	return http.ListenAndServe(s.config.WebPort, s.mux)

}

// ConfigureDB ...
func (s *Server) ConfigureDB() error {
	db := database.NewDB(s.config.Database)
	if err := db.InitDB(); err != nil {
		return err
	}
	s.database = db //fill Server with DB instance
	return nil
}
