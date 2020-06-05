package controller

import (
	"net/http"

	"github.com/astgot/forum/internal/model"
)

// MainHandle ...
func (m *Multiplexer) MainHandle() http.HandlerFunc {

	// Need to create structure to show array of Users, Posts for arranging them in HTML
	type PostRaw struct {
		Post *model.Post
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

			guest := &PostRaw{}
			for _, post := range posts {
				guest.Post = m.db.GetPostByPID(post.PostID)
				mainPage.PostScroll = append(mainPage.PostScroll, guest)
			}
			// if err := tpl.ExecuteTemplate(w, "main.html", mainPage); err != nil {
			// 	http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
			// 	return
			// }
			tpl.ExecuteTemplate(w, "main.html", mainPage)
			return

		}
		user, _ := m.db.GetUserByCookie(cookie.Value)
		mainPage.AuthUser = user
		auth := &PostRaw{}
		for _, post := range posts {
			auth.Post = m.db.GetPostByPID(post.PostID)
			mainPage.PostScroll = append(mainPage.PostScroll, auth)

		}
		// if err := tpl.ExecuteTemplate(w, "main.html", mainPage); err != nil {
		// 	http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
		// 	return
		// }
		tpl.ExecuteTemplate(w, "main.html", mainPage)

	}
}
