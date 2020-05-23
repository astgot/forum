package server

import (
	"io"
	"net/http"

	"github.com/astgot/forum/internal/database"
	"github.com/astgot/forum/internal/model"
)

// ConfigureRouter ...
func (s *Server) ConfigureRouter() {

	s.mux.HandleFunc("/", s.MainHandle())
	s.mux.HandleFunc("/signup", s.SignupHandle())
	s.mux.HandleFunc("/confirmation", ConfirmHandler)
	return
}

// MainHandle ...
func (s *Server) MainHandle() http.HandlerFunc {

	// Here we can create our own struct, which is usable only here
	return func(w http.ResponseWriter, r *http.Request) {
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

			encryptPass := model.HashPwd(userInfo.Password)
			userInfo.EncryptedPwd = encryptPass // fill with Encrypted Password
			userRepo := database.UsersStore{}
			err := userRepo.Create(&userInfo)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotImplemented)
				return
			}

			http.Redirect(w, r, "/confirmation", http.StatusCreated)

		}
	}

}

// LoginHandle --->
func (s *Server) LoginHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/login" {
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}

		switch r.Method {
		case "GET":
			tpl.ExecuteTemplate(w, "login.html", nil)
		case "POST":
			r.ParseForm()
			login := model.Users{
				Username: r.PostFormValue("Username"),
				Email:    r.PostFormValue("Email"),
				Password: r.PostFormValue("Password"),
			}
			find := database.UsersStore{}
			if login.Username == "" {
				_, err := find.FindByEmail(login.Email)
				if err != nil {
					http.Error(w, "Invalid", http.StatusUnauthorized)
					return
				}
			} else if login.Email == "" {
				_, err := find.FindByUsername(login.Username)
				if err != nil {
					http.Error(w, "Invalid", http.StatusUnauthorized)
					return
				}
			}

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
