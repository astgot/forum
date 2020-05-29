package controller

import (
	"net/http"
	"time"

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
		Name:    sessionName,
		Value:   GenerateSessionToken(),
		Expires: time.Now().Add(20 * time.Minute),
	}

	http.SetCookie(w, cookieSession)
	if sessionName != "guest" {
		m.db.InsertSession(user, cookieSession)
	}

}
