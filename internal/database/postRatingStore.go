package database

import (
	"fmt"

	"github.com/astgot/forum/internal/model"
)

// AddRateToPost ...
func (d *Database) AddRateToPost(l *model.PostRating, uid int64) bool {
	/*Need to handle the situation:
	If user liked the post, DB will insert record with likeCount value only.
	And then for example, another user will dislike the post, DB will add another record
	with this postID and likeCount as null
	*/

	rate := d.GetRateCountOfPost(l.PostID)
	fmt.Println(rate.DislikeCount, rate.LikeCount, "rates")
	var kind int64

	if l.Like {
		kind = 1
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
			if ok := d.AddLike(rate.LikeCount, rate.PostID); !ok {
				return false
			}
		}
		// If dislike
	} else {
		kind = 0
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
			if ok := d.AddDislike(rate.DislikeCount, rate.PostID); !ok {
				return false
			}
		}

	}

	stmnt, err := d.db.Prepare("INSERT INTO RateUserPost (userID, postID, kind) VALUES (?, ?)")
	_, err = stmnt.Exec(uid, l.PostID, kind)
	if err != nil {
		fmt.Println("RateUserPost error")
		return false
	}

	return true
}

// DeleteRateFromPost ...
func (d *Database) DeleteRateFromPost(rate *model.PostRating, uid int64) bool {
	// 1) What user did now? (like or dislike)
	// 2) What user have done before?
	// if user 1) liked and 2) liked ---> Delete like from post
	// If 1) liked 2) dislike ---> Delete like and add dislike
	// If 1) disliked 2)disliked ---> Delete dislike
	// If 1) disliked 2) liked ---> Delete dislike and add like

	var before int64
	if err := d.db.QueryRow("SELECT kind FROM RateUserPost WHERE userID=? AND postID=?", uid, rate.PostID).
		Scan(&before); err != nil {
		fmt.Println("DeleteRateFromPost error type")
		return false
	}
	rateCount := d.GetRateCountOfPost(rate.PostID)
	rate.DislikeCount = rateCount.DislikeCount
	rate.LikeCount = rateCount.LikeCount

	// Scenarios
	if before == 1 && rate.Like {
		// delete like
		if ok := d.DeleteLike(rate.LikeCount, rate.PostID); !ok {
			return false
		}

	} else if before == 1 && !(rate.Like) {
		//delete like, add dislike
		if ok := d.DeleteLike(rate.LikeCount, rate.PostID); !ok {
			return false
		}
		if ok := d.AddDislike(rate.DislikeCount, rate.PostID); !ok {
			return false
		}

	} else if before == 0 && !(rate.Like) {
		//delete dislike
		if ok := d.DeleteDislike(rate.DislikeCount, rate.PostID); !ok {
			return false
		}
	} else if before == 0 && rate.Like {
		//delete dislike, add like
		if ok := d.DeleteDislike(rate.DislikeCount, rate.PostID); !ok {
			return false
		}
		if ok := d.AddLike(rate.LikeCount, rate.PostID); !ok {
			return false
		}
	}

	return true
}
