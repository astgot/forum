package database

import "log"

// BuildSchema ---> create Tables
func (d *Database) BuildSchema() error {
	users, err := d.db.Prepare(`CREATE TABLE IF NOT EXISTS Users (
									id INTEGER PRIMARY KEY NOT NULL, 
									firstname TEXT NOT NULL, 
									lastname TEXT NOT NULL, 
									username TEXT NOT NULL UNIQUE, 
									email TEXT NOT NULL UNIQUE, 
									password TEXT NOT NULL
								)`)

	defer users.Close()
	CheckErr(err)
	users.Exec()

	sessions, err := d.db.Prepare(`CREATE TABLE IF NOT EXISTS Sessions (
		userID INTEGER NOT NULL,
		cookieName TEXT NOT NULL,
		cookieValue TEXT NOT NULL UNIQUE,
		FOREIGN KEY(userID) REFERENCES Users(id)
	)`)
	defer sessions.Close()
	CheckErr(err)
	sessions.Exec()

	return nil
}

// CheckErr ...
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
