package database

import "fmt"

// InsertPostMapInfo ...
func (d *Database) InsertPostMapInfo(postID, threadID int64) error {
	stmnt, err := d.db.Prepare("INSERT INTO PostMapping (postID, threadID) VALUES (?, ?)")
	defer stmnt.Close()
	if err != nil {
		fmt.Println("postmap insert error")
		return err
	}
	stmnt.Exec(postID, threadID)
	return nil
}
