package controller

import (
	"net/http"

	"github.com/astgot/forum/internal/model"
)

// MainHandle ...
func (m *Multiplexer) MainHandle() http.HandlerFunc {

	// Need to create structure to show array of Users, Posts, Comments, Categories for arranging them in HTML
	type PostRaw struct {
		Post     *model.Post
		Threads  []*model.Thread
		PostRate *model.PostRating
		Comments []*model.Comments
	}
	var mainPage struct {
		AuthUser   *model.Users
		PostScroll []*PostRaw
	}
	// Here we can create our own struct, which is usable only here
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "/main" {
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}
		posts := m.GetAllPosts(w)

		cookie, err := r.Cookie("authenticated")
		if err != nil {
			// if user is guest, retrieve all posts for displaying
			for _, post := range posts {
				guest := &PostRaw{}
				guest.Comments, _ = m.db.GetCommentsOfPost(post.PostID)
				// guest.Author, _ = m.db.FindByUserID(post.UserID)
				guest.Post = post
				guest.Threads, _ = m.db.GetThreadOfPost(post.PostID)
				guest.PostRate = m.db.GetRateCountOfPost(post.PostID)
				mainPage.PostScroll = append(mainPage.PostScroll, guest)
			}
			tpl.ExecuteTemplate(w, "main.html", mainPage)
			return

		}
		// if User is authenticated
		user, _ := m.db.GetUserByCookie(cookie.Value)
		mainPage.AuthUser = user
		auth := &PostRaw{}
		for _, post := range posts {
			auth.Post, _ = m.db.GetPostByPID(post.PostID)
			auth.Threads, err = m.db.GetThreadOfPost(post.PostID)
			if err != nil {
				http.Error(w, "Threads retrieving error", http.StatusInternalServerError)
				return
			}
			mainPage.PostScroll = append(mainPage.PostScroll, auth)

		}

		tpl.ExecuteTemplate(w, "main.html", mainPage)

	}
}
