package controller

import (
	"html/template"
	"net/http"

	"github.com/astgot/forum/internal/database"
)

var tpl = template.Must(template.ParseGlob("web/templates/*"))

// Warning ...
var Warning struct {
	Warn string
}

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

// CreateHandlers ...
func (m *Multiplexer) CreateHandlers() {
	fs := http.FileServer(http.Dir("web/css"))
	m.Mux.Handle("/css/", http.StripPrefix("/css/", fs))
	m.Mux.HandleFunc("/", m.MainHandle())
	m.Mux.HandleFunc("/signup", m.SignupHandle())
	m.Mux.HandleFunc("/login", m.LoginHandle())
	m.Mux.HandleFunc("/logout", m.LogoutHandle())
	m.Mux.HandleFunc("/confirmation", ConfirmHandler)
	m.Mux.HandleFunc("/profile", m.ProfileHandler())
	m.Mux.HandleFunc("/create", m.CreatePostHandler())
	m.Mux.HandleFunc("/post", m.PostView())
	m.Mux.HandleFunc("/rate", m.RateHandler())
	m.Mux.HandleFunc("/filter", m.FilterHandler())
}

// ConfigureDB ...
func (m *Multiplexer) ConfigureDB() error {
	if err := m.db.InitDB(); err != nil {
		return err
	}
	return nil
}

// WarnMessage ...
func WarnMessage(w http.ResponseWriter, warn string) {
	Warning.Warn = warn
	tpl.ExecuteTemplate(w, "error.html", Warning)
}
