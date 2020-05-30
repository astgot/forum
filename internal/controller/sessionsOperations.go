package controller

import (
	"fmt"
	"net/http"

	"github.com/astgot/forum/internal/model"
	uuid "github.com/satori/go.uuid"
)

// GenerateSessionToken ...
func GenerateSessionToken() string {
	return uuid.NewV4().String()
}

// CheckSession ...
func CheckSession(r *http.Request, sessionName string) error {
	_, err := r.Cookie(sessionName)
	if err != nil {
		return err
	}
	return nil
}

// AddSession ...
func (m *Multiplexer) AddSession(w http.ResponseWriter, sessionName string, user *model.Users) {
	cookieSession := &http.Cookie{
		Name:   sessionName,
		Value:  GenerateSessionToken(),
		MaxAge: 900,
	}

	http.SetCookie(w, cookieSession)
	if sessionName != "guest" {
		if _, err := m.db.InsertSession(user, cookieSession); err != nil {
			fmt.Println("Error")
		}
	}

}

// DeleteSession ...
func (m *Multiplexer) DeleteSession(w http.ResponseWriter, sessionValue string) {
	cookie := &http.Cookie{
		Name:   "null",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	m.db.DeleteCookieFromDB(sessionValue)
}
