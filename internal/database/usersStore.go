package database

import (
	"fmt"
	"net/http"

	"github.com/astgot/forum/internal/model"
)

// Create ---> signing up
func (d *Database) Create(u *model.Users) (*model.Users, error) {
	if err := d.Open(); err != nil {
		return nil, err
	}

	stmnt, err := d.db.Prepare("INSERT INTO Users (firstname, lastname, username, email, password) VALUES (?, ?, ?, ?, ?)")
	res, err := stmnt.Exec(u.Firstname, u.Lastname, u.Username, u.Email, u.EncryptedPwd)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	id, _ := res.LastInsertId()
	u.ID = id

	return u, nil

}

// FindByUsername ...
func (d *Database) FindByUsername(username string) (*model.Users, error) {
	if err := d.Open(); err != nil {
		return nil, err
	}

	u := &model.Users{}
	if err := d.db.QueryRow("SELECT id, username, password FROM Users where username = ?", username).Scan(
		&u.ID,
		&u.Username,
		&u.EncryptedPwd,
	); err != nil {
		return nil, err
	}

	return u, nil
}

// FindByEmail ...
func (d *Database) FindByEmail(email string) (*model.Users, error) {
	if err := d.Open(); err != nil {
		return nil, err
	}
	u := &model.Users{}
	if err := d.db.QueryRow("SELECT id, email, password FROM Users where email = ?", email).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPwd,
	); err != nil {
		fmt.Println("asas")
		return nil, err
	}

	return u, nil
}

// GetUserID ...
func (d *Database) GetUserID(user *model.Users, email bool) int64 {
	if err := d.Open(); err != nil {
		return 0
	}

	if email {
		if err := d.db.QueryRow("SELECT id FROM Users where email = ?", user.Username).Scan(
			&user.ID,
		); err != nil {
			return 0
		}
	} else {
		if err := d.db.QueryRow("SELECT id FROM Users where username = ?", user.Username).Scan(
			&user.ID,
		); err != nil {
			return 0
		}

	}
	return user.ID

}

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
