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

// AddRateToPost ...
func (d *Database) AddRateToPost(l *model.PostRating, uid int64) bool {
	/*Need to handle the situation:
	If user liked the post, DB will insert record with likeCount value only.
	And then for example, another user will dislike the post, DB will add another record
	with this postID and likeCount as null
	*/

	rate := d.GetRateCountOfPost(l.PostID)
	fmt.Println(rate.DislikeCount, rate.LikeCount, "rates")

	if l.Like {
		if rate.LikeCount == 0 && rate.DislikeCount == 0 {
			stmnt, err := d.db.Prepare("INSERT INTO PostRating (postID, likeCount, dislikeCount) VALUES (?, ?, ?)")
			defer stmnt.Close()
			_, err = stmnt.Exec(l.PostID, rate.LikeCount+1, rate.DislikeCount)
			if err != nil {
				fmt.Println("db Insert PostRating error", err.Error())
				return false
			}

		} else {
			// Update column "likeCount" in the table
			stmnt, err := d.db.Prepare("UPDATE PostRating SET likeCount=? WHERE postID=?")
			defer stmnt.Close()
			_, err = stmnt.Exec(rate.LikeCount+1, l.PostID)
			if err != nil {
				fmt.Println("update likecount error")
				return false
			}
		}
		// If dislike
	} else {
		if rate.LikeCount == 0 && rate.DislikeCount == 0 {
			stmnt, err := d.db.Prepare("INSERT INTO PostRating (postID, likeCount, dislikeCount) VALUES (?, ?, ?)")
			defer stmnt.Close()
			_, err = stmnt.Exec(l.PostID, rate.LikeCount, rate.DislikeCount+1)
			if err != nil {
				fmt.Println("db Insert PostRating dislike error")
				return false
			}
		} else {
			// Update column "dislikeCount" in the table
			stmnt, err := d.db.Prepare("UPDATE PostRating SET dislikeCount=? WHERE postID=?")
			defer stmnt.Close()
			_, err = stmnt.Exec(rate.DislikeCount+1, l.PostID)
			if err != nil {
				fmt.Println("update dislikecount error")
				return false
			}
		}

	}

	stmnt, err := d.db.Prepare("INSERT INTO RateUserPost (userID, postID) VALUES (?, ?)")
	_, err = stmnt.Exec(uid, l.PostID)
	if err != nil {
		fmt.Println("RateUserPost error")
		return false
	}

	return true
}

// AddRateToComment ...
func (d *Database) AddRateToComment(l *model.CommentRating) *model.CommentRating {
	rate := model.NewCommentRating()
	if l.Like {

	} else {

	}
	return rate
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

// IsUserRateComm ...
func (d *Database) IsUserRateComm(uid, pid, cid int64) bool {
	return false
}
