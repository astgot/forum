package tools

import (
	"net/http"
	"time"

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
func AddSession(w http.ResponseWriter, sessionName string) {
	cookie := &http.Cookie{
		Name:    sessionName,
		Value:   GenerateSessionToken(),
		Expires: time.Now().Add(20 * time.Minute),
	}

	http.SetCookie(w, cookie)

}
