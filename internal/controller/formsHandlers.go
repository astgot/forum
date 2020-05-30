package controller

import (
	"html/template"
	"net/http"

	"github.com/astgot/forum/internal/model"
)

var tpl = template.Must(template.ParseGlob("web/templates/*"))

// MainHandle ...
func (m *Multiplexer) MainHandle() http.HandlerFunc {

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
		if err := CheckSession(r, "guest"); err == http.ErrNoCookie {
			m.AddSession(w, "guest", nil)
		}

	}
}

// SignupHandle ---> /signup
func (m *Multiplexer) SignupHandle() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/signup" {
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}
		if r.Method == "GET" {
			tpl.ExecuteTemplate(w, "signup.html", nil)
		} else if r.Method == "POST" {
			r.ParseForm() // Parsing Form from the front-end
			user := &model.Users{
				Firstname:  r.PostFormValue("Firstname"),
				Lastname:   r.PostFormValue("Lastname"),
				Username:   r.PostFormValue("Username"),
				Email:      r.PostFormValue("Email"),
				Password:   r.PostFormValue("Password"),
				ConfirmPwd: r.PostFormValue("Confirm"),
			}

			if ValidateInput(user) == false {
				tpl.ExecuteTemplate(w, "signup.html", user)
				return
			}

			encryptPass := HashPassword(user.Password)
			user.EncryptedPwd = encryptPass   // fill with Encrypted Password
			newUser, err := m.db.Create(user) // Sending
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			m.AddSession(w, "guest", newUser) // guest session
			http.Redirect(w, r, "/confirmation", http.StatusSeeOther)

		}
	}

}

// LoginHandle ---> /login
func (m *Multiplexer) LoginHandle() http.HandlerFunc {
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
			login := &model.Users{
				Username: r.PostFormValue("Username"),
				Password: r.PostFormValue("Password"),
			}
			check := Login{}

			check.unameOrEmail = UnameOrEmail(login.Username)

			if check.unameOrEmail {
				u, err := m.db.FindByEmail(login.Username)
				if err != nil {
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
				check.auth = ComparePassword(u.EncryptedPwd, login.Password)
				if !check.auth {
					http.Error(w, "Invalid credentials", http.StatusUnauthorized)
					return
				}
			} else if !check.unameOrEmail {
				u, err := m.db.FindByUsername(login.Username)
				if err != nil {
					http.Error(w, "Invalid credentials", http.StatusUnauthorized)
					return
				}
				check.auth = ComparePassword(u.EncryptedPwd, login.Password)
				if !check.auth {
					http.Error(w, "Invalid credentials", http.StatusUnauthorized)
					return
				}

			}
			login.ID = m.db.GetUserID(login, check.unameOrEmail)
			m.AddSession(w, login.Username, login) // Add cookie session after successful authentication
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
