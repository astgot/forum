package controller

import (
	"net/http"

	"github.com/astgot/forum/internal/model"
)

// MainHandle ...
func (m *Multiplexer) MainHandle() http.HandlerFunc {
	type Postview struct {
		Posts []model.Post
	}
	// Here we can create our own struct, which is usable only here
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "/main" {
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}
		if err := tpl.ExecuteTemplate(w, "main.html", nil); err != nil {
			http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
			return
		}
		err := m.CheckSession(r, "guest")
		if err == http.ErrNoCookie {
			m.AddSession(w, "guest", nil)
			return
		}

	}
}
