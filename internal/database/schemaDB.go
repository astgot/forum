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
	CheckErr(err)
	users.Exec()
	// d.Close()
	return nil
}

// CheckErr ...
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
