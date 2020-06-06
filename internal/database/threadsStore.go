package database

import (
	"fmt"

	"github.com/astgot/forum/internal/model"
)

// InsertThreadInfo ..
func (d *Database) InsertThreadInfo(thread *model.Thread, postID int64) error {
	stmnt, err := d.db.Prepare("INSERT INTO Threads (Name) VALUES (?)")
	if err != nil {
		fmt.Println("insert Threads error")
		return err
	}
	defer stmnt.Close()
	res, err := stmnt.Exec(thread.Name)
	if err != nil {
		fmt.Println(err.Error(), "---> exec Threads error")
		if err.Error() == "UNIQUE constraint failed: Threads.Name" {
			ID := d.GetThreadID(thread.Name)
			d.InsertPostMapInfo(postID, ID)
		}
		return err
	}
	threadID, _ := res.LastInsertId()
	d.InsertPostMapInfo(postID, threadID)

	return nil
}

// GetAllThreads ...
func (d *Database) GetAllThreads() []*model.Thread {
	var threads []*model.Thread
	res, err := d.db.Query("SELECT * FROM Threads")
	if err != nil {
		fmt.Println("thread query error")
		return nil
	}
	for res.Next() {
		thread := model.NewThread()
		if err := res.Scan(&thread.ID, &thread.Name); err != nil {
			fmt.Println("thread scan error")
			return nil
		}
		threads = append(threads, thread)

	}
	return threads
}

// GetThreadID ...
func (d *Database) GetThreadID(name string) int64 {
	var id int64
	if err := d.db.QueryRow("SELECT ID FROM Threads WHERE Name = ?", name).Scan(&id); err != nil {
		fmt.Println("threadID retrieve error")
	}

	return id
}
