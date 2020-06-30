package database

import (
	"fmt"

	"github.com/astgot/forum/internal/model"
)

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

// GetThreadOfPost ...
func (d *Database) GetThreadOfPost(postID int64) ([]*model.Thread, error) {
	threadIDs := []int64{}      // for getting all threadIDs from postID
	var threads []*model.Thread // for getting names of threads
	res, err := d.db.Query("SELECT threadID FROM PostMapping WHERE postID = ?", postID)
	if err != nil { // if err == sql.ErrNoRows ---> if no category in the post
		return nil, err
	}
	defer res.Close()
	/* Here we retrieve all threads relating with one single post*/
	for res.Next() {
		postMap := model.NewPostMap()
		if err := res.Scan(&postMap.ThreadID); err != nil {
			fmt.Println("error func\"GetThreadOfPost()\"")
			return nil, err
		}
		threadIDs = append(threadIDs, postMap.ThreadID)
	}
	// After getting IDs of threads, by these IDs we will get names of the threads
	thread := model.NewThread()
	for _, threadID := range threadIDs {
		thread, err = d.GetThreadByID(threadID) // exactly here
		if err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}

	return threads, nil
}
