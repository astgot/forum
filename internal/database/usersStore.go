package database

import (
	"fmt"

	"github.com/astgot/forum/internal/model"
)

// UsersStore will store users' data
type UsersStore struct {
	database *Database
}

// Create ---> signing up
func (us *UsersStore) Create(u *model.Users) error {
	fmt.Println(u.Firstname, u.Lastname, u.Username, u.Email, u.EncryptedPwd)
	res, err := us.database.db.Exec("INSERT INTO Users (firstname, lastname, username, email, password) VALUES (?, ?, ?, ?, ?)",
		u.Firstname, u.Lastname, u.Username, u.Email, u.EncryptedPwd)
	if err != nil {
		return err
	}
	fmt.Println("assssas")

	id, _ := res.LastInsertId()

	u.ID = int(id)
	// return u, nil
	return nil
}

// FindByUsername ...
func (us *UsersStore) FindByUsername(username string) (*model.Users, error) {

	u := &model.Users{}
	if err := us.database.db.QueryRow("SELECT id, username, password FROM Users where username = \"?\"", username).Scan(
		&u.ID,
		&u.Username,
		&u.Password,
	); err != nil {
		return nil, err
	}

	return u, nil
}

// FindByEmail ...
func (us *UsersStore) FindByEmail(email string) (*model.Users, error) {
	u := &model.Users{}
	if err := us.database.db.QueryRow("SELECT id, email, password FROM Users where email = \"?\"", email).Scan(
		&u.ID,
		&u.Email,
		&u.Password,
	); err != nil {
		return nil, err
	}

	return u, nil
}
