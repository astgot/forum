package database

import (
	"fmt"

	"github.com/astgot/forum/internal/model"
)

// Create ---> signing up
func (d *Database) Create(u *model.Users) error {
	if err := d.Open(); err != nil {
		return err
	}

	stmnt, err := d.db.Prepare("INSERT INTO Users (firstname, lastname, username, email, password) VALUES (?, ?, ?, ?, ?)")
	res, err := stmnt.Exec(u.Firstname, u.Lastname, u.Username, u.Email, u.EncryptedPwd)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	id, _ := res.LastInsertId()
	u.ID = id

	// return u, nil
	return nil
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
