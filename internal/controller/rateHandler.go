package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/astgot/forum/internal/model"
)

// LikeHandler ...
func (m *Multiplexer) LikeHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rate" {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		var Islike = true
		var isPost = true
		cookie, err := r.Cookie("authenticated")
		if err != nil {
			http.Error(w, "You need to authorize", http.StatusForbidden)
			return
			// maybe make template showing to authorize firstly
		}

		id, err := strconv.Atoi(r.URL.Query().Get("post_id"))
		if err != nil {
			id, err = strconv.Atoi(r.URL.Query().Get("comment_id"))
			if err != nil {
				http.Error(w, "Invalid parameter", http.StatusBadRequest)
				return
			}
			isPost = false

		}
		// if "id" is negative number, it will be dislike
		if id < 0 {
			Islike = false
			id *= -1
		}
		user, err := m.db.GetUserByCookie(cookie.Value)
		if err != nil {
			fmt.Println("GetUserByCookie rateHandler.go error")
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		if isPost {
			like := model.NewPostRating()
			like.Like = Islike
			_, err := m.db.GetPostByPID(int64(id))
			if err != nil {
				http.Error(w, "The post not found", http.StatusBadRequest)
				fmt.Println("GetPostByPID error")
				return
			}
			like.PostID = int64(id)
			// like.UID = user.ID // assign like to UserID --> to check user liked this post,
			// prevent multiple liking of the post
			// Need to return new rate count
			/* Check user liked this post
			if Yes, delete rate from the post
			*/
			isRated := m.db.IsUserRatePost(user.ID, int64(id))
			if isRated {
				if ok := m.db.DeleteRateFromPost(like, user.ID); !ok {
					fmt.Println("DeleteRateOfPost error")
					http.Error(w, "Something went wrong", http.StatusInternalServerError)
					return
				}
			} else {
				if ok := m.db.AddRateToPost(like, user.ID); !ok {
					http.Error(w, "Something went wrong", http.StatusInternalServerError)
					fmt.Println("AddRatePost() error")
					return
				}
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)

		} else {
			like := model.NewCommentRating()
			like.Like = Islike
			comment, err := m.db.GetCommentByID(int64(id))
			if err != nil {
				http.Error(w, "The comment not found", http.StatusBadRequest)
				return
			}
			like.CommentID = comment.CommentID
			// Need to retrieve post from comment id

			m.db.AddRateToComment(like) // Need to return new rate count
		}

	}
}
