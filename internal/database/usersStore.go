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
// func (d *Database) GetUserID(creds string, uname bool) int64 {
// 	if err := d.Open(); err != nil {
// 		return 0
// 	}
// 	var id int64
// 	if uname {
// 		if err := d.db.QueryRow("SELECT id FROM Users where email = ?", creds).Scan(
// 			id,
// 		); err != nil {
// 			return 0
// 		}
// 	} else {
// 		if err := d.db.QueryRow("SELECT id FROM Users where username = ?", creds).Scan(
// 			id,
// 		); err != nil {
// 			return 0
// 		}

// 	}
// 	return id

// }

// InsertSession ...
func (d *Database) InsertSession(u *model.Users, session *http.Cookie) (*model.Sessions, error) {
	if err := d.Open(); err != nil {
		return nil, err
	}
	cookie := model.NewSession()
	if err := d.db.QueryRow("INSERT INTO Sessions (cookieName, cookieValue) VALUES ( ?, ?)", session.Name, session.Value).Scan(
		&cookie.SessionName,
		&cookie.SessionValue,
	); err != nil {
		return nil, err
	}
	return cookie, nil
}
