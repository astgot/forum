package server

import (
	"io"
	"net/http"

	"github.com/astgot/forum/internal/model"
)

// ConfigureRouter ...
func (s *Server) ConfigureRouter() {

	s.mux.HandleFunc("/", s.MainHandle())
	s.mux.HandleFunc("/signup", s.SignupHandle())
	s.mux.Handle("/login", s.LoginHandle())
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
		io.WriteString(w, "<b>YO Wazzup!!!</b>")
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
			userInfo := model.Users{
				Firstname:  r.PostFormValue("Firstname"),
				Lastname:   r.PostFormValue("Lastname"),
				Username:   r.PostFormValue("Username"),
				Email:      r.PostFormValue("Email"),
				Password:   r.PostFormValue("Password"),
				ConfirmPwd: r.PostFormValue("Confirm"),
			}

			if userInfo.Validate() == false {
				tpl.ExecuteTemplate(w, "signup.html", userInfo)
				return
			}

			encryptPass := model.HashPassword(userInfo.Password)
			userInfo.EncryptedPwd = encryptPass // fill with Encrypted Password
			err := s.database.User().Create(&userInfo)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			http.Redirect(w, r, "/confirmation", http.StatusSeeOther)

		}
	}

}

// LoginHandle --->
func (s *Server) LoginHandle() http.HandlerFunc {
	type Login struct {
		auth         bool
		unameOremail bool
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

			check.unameOremail = model.UnameOrEmail(login.Username)

			if check.unameOremail {
				u, err := s.database.User().FindByEmail(login.Username)
				if err != nil {
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
				check.auth = model.ComparePassword(u.EncryptedPwd, login.Password)
				if !check.auth {
					http.Error(w, "Invalid credentials", http.StatusUnauthorized)
					return
				}
			} else if !check.unameOremail {
				u, err := s.database.User().FindByUsername(login.Username)
				if err != nil {
					http.Error(w, "Invalid credentials", http.StatusUnauthorized)
					return
				}
				check.auth = model.ComparePassword(u.EncryptedPwd, login.Password)
				if !check.auth {
					http.Error(w, "Invalid credentials", http.StatusUnauthorized)
					return
				}

			}
			http.Redirect(w, r, "/main", http.StatusSeeOther)

		}

	}
}

// ConfirmHandler -->
func ConfirmHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/confirmation" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	tpl.ExecuteTemplate(w, "confirmation.html", nil)
}
