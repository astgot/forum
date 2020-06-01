package database

import "github.com/astgot/forum/internal/model"

// InsertPostInfo ...
func (d *Database) InsertPostInfo(p *model.Post) error {
	if err := d.Open(); err != nil {
		return err
	}

	stmnt, err := d.db.Prepare("INSERT INTO PostsInfo (userID, creationDate) VALUES (?, ?)")
	defer stmnt.Close()
	if err != nil {
		return err
	}
	stmnt.Exec(p.UserID, p.CreationDate)
	return nil
}
