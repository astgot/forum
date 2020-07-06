package database

import (
	"fmt"

	"github.com/astgot/forum/internal/model"
)

// GetRateCountOfPost ...
func (d *Database) GetRateCountOfPost(postID int64) *model.PostRating {
	// Здесь нужно разобраться
	rates := model.NewPostRating()
	if err := d.db.QueryRow("SELECT * FROM PostRating WHERE postID = ?", postID).
		Scan(&rates.PostID,
			&rates.LikeCount,
			&rates.DislikeCount,
		); err != nil {
		// It means nobody rated the post, likeCount and dislikeCount now is zero
		rates.LikeCount = 0
		rates.DislikeCount = 0
	}
	return rates
}

// IsUserRatePost ...
func (d *Database) IsUserRatePost(uid, pid int64) bool {
	var comp int64
	// need to check all rated posts of the user
	res, err := d.db.Query("SELECT postID FROM RateUserPost WHERE userID=?", uid, pid)
	if err != nil {
		return false
	}
	defer res.Close()
	// Check postID is rated or not
	for res.Next() {
		if err := res.Scan(&comp); err != nil {
			return false
		}
		if comp == pid {
			return true
		}
		comp = 0
	}
	return false
}

// AddLike ...
func (d *Database) AddLike(likeCnt, postID int64) bool {
	stmnt, err := d.db.Prepare("UPDATE PostRating SET likeCount=? WHERE postID=?")
	defer stmnt.Close()
	_, err = stmnt.Exec(likeCnt+1, postID)
	if err != nil {
		fmt.Println("update likecount error")
		return false
	}
	return true
}

// DeleteLike ...
func (d *Database) DeleteLike(likeCnt, postID int64) bool {
	stmnt, err := d.db.Prepare("UPDATE PostRating SET likeCount=? WHERE postID=?")
	defer stmnt.Close()
	_, err = stmnt.Exec(likeCnt-1, postID)
	if err != nil {
		fmt.Println("update likecount error")
		return false
	}
	return true
}

// AddDislike ...
func (d *Database) AddDislike(dislikeCnt, postID int64) bool {
	stmnt, err := d.db.Prepare("UPDATE PostRating SET dislikeCount=? WHERE postID=?")
	defer stmnt.Close()
	_, err = stmnt.Exec(dislikeCnt+1, postID)
	if err != nil {
		fmt.Println("update dislikecount error")
		return false
	}
	return true
}

// DeleteDislike ...
func (d *Database) DeleteDislike(dislikeCnt, postID int64) bool {
	stmnt, err := d.db.Prepare("UPDATE PostRating SET dislikeCount=? WHERE postID=?")
	defer stmnt.Close()
	_, err = stmnt.Exec(dislikeCnt-1, postID)
	if err != nil {
		fmt.Println("update dislikecount error")
		return false
	}
	return true
}
