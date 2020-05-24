package server

import (
	"net/http"

	"github.com/astgot/forum/internal/model"
	"github.com/astgot/forum/internal/tools"
)

// ConfigureRouter ...
func (s *Server) ConfigureRouter() {

	s.mux.HandleFunc("/", s.MainHandle())
	s.mux.HandleFunc("/signup", s.SignupHandle())
	s.mux.HandleFunc("/login", s.LoginHandle())
	s.mux.HandleFunc("/confirmation", ConfirmHandler)
	return
}

// MainHandle ...
func (s *Server) MainHandle() http.HandlerFunc {

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
		// Checking for session, processing ...
		if err := tools.CheckSession(r, "authenticated"); err != nil {
			if err = tools.CheckSession(r, "guest"); err != nil {
				tools.AddSession(w, "guest")
			}
		}

	}
}

// SignupHandle ---> /signup
func (s *Server) SignupHandle() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/signup" {
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}
		if r.Method == "GET" {
			tpl.ExecuteTemplate(w, "signup.html", nil)
		} else if r.Method == "POST" {
			r.ParseForm() // Parsing Form from the front-end
			userInfo := &model.Users{
				Firstname:  r.PostFormValue("Firstname"),
				Lastname:   r.PostFormValue("Lastname"),
				Username:   r.PostFormValue("Username"),
				Email:      r.PostFormValue("Email"),
				Password:   r.PostFormValue("Password"),
				ConfirmPwd: r.PostFormValue("Confirm"),
			}

			if tools.ValidateInput(userInfo) == false {
				tpl.ExecuteTemplate(w, "signup.html", userInfo)
				return
			}

			encryptPass := tools.HashPassword(userInfo.Password)
			userInfo.EncryptedPwd = encryptPass       // fill with Encrypted Password
			err := s.database.User().Create(userInfo) // Sending
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			tools.AddSession(w, "guest") // guest session
			http.Redirect(w, r, "/confirmation", http.StatusSeeOther)

		}
	}

}

// LoginHandle ---> /login
func (s *Server) LoginHandle() http.HandlerFunc {
	type Login struct {
		auth         bool
		unameOrEmail bool
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/login" {
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}

		if r.Method == "GET" {
			tpl.ExecuteTemplate(w, "login.html", nil)
		} else if r.Method == "POST" {
			r.ParseForm()
			login := model.Users{
				Username: r.PostFormValue("Username"),
				Password: r.PostFormValue("Password"),
			}
			check := Login{}

			check.unameOrEmail = tools.UnameOrEmail(login.Username)

			if check.unameOrEmail {
				u, err := s.database.User().FindByEmail(login.Username)
				if err != nil {
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
				check.auth = tools.ComparePassword(u.EncryptedPwd, login.Password)
				if !check.auth {
					http.Error(w, "Invalid credentials", http.StatusUnauthorized)
					return
				}
			} else if !check.unameOrEmail {
				u, err := s.database.User().FindByUsername(login.Username)
				if err != nil {
					http.Error(w, "Invalid credentials", http.StatusUnauthorized)
					return
				}
				check.auth = tools.ComparePassword(u.EncryptedPwd, login.Password)
				if !check.auth {
					http.Error(w, "Invalid credentials", http.StatusUnauthorized)
					return
				}

			}
			tools.AddSession(w, "authenticated") // Add cookie session after successful authentication
			http.Redirect(w, r, "/main", http.StatusSeeOther)

		}

	}
}

// ConfirmHandler --> /confirmation
func ConfirmHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/confirmation" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	tpl.ExecuteTemplate(w, "confirmation.html", nil)
}
