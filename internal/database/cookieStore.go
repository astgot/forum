package database

import (
	"net/http"

	"github.com/astgot/forum/internal/model"
)

// InsertSession ...
func (d *Database) InsertSession(u *model.Users, session *http.Cookie) (*model.Sessions, error) {
	if err := d.Open(); err != nil {
		return nil, err
	}
	cookie := model.NewSession()
	if err := d.db.QueryRow("INSERT INTO Sessions (userID, cookieName, cookieValue) VALUES (?, ?, ?)", u.ID, session.Name, session.Value).Scan(
		&cookie.UserID,
		&cookie.SessionName,
		&cookie.SessionValue,
	); err != nil {
		return nil, err
	}
	return cookie, nil
}

// DeleteCookieFromDB ...
func (d *Database) DeleteCookieFromDB(cookieValue string) error {
	stmnt, err := d.db.Prepare("DELETE FROM Sessions WHERE cookieValue = ?")
	stmnt.Exec(cookieValue)
	if err != nil {
		return err
	}
	return nil

}
