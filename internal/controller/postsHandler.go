package controller

import (
	"net/http"
	"time"

	"github.com/astgot/forum/internal/model"
)

// CreatePostHandler ...
func (m *Multiplexer) CreatePostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/create" {
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}
		cookie, err := r.Cookie("authenticated")
		if err != nil {
			http.Error(w, "Not Authorized", http.StatusUnauthorized)
			return
			// OR tpl.ExecuteTemplate(w, "error.html", nil) // need to add error "Need to make authorization"
		}

		u, err := m.db.GetUserByCookie(cookie.Value)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError) // Check workflow of DB
			return
		}
		post := model.NewPost()
		if r.Method == "POST" {
			r.ParseForm()
			post.UserID = u.ID
			post.Author = u.Firstname + " " + u.Lastname + " aka " + "\"" + u.Username + "\""
			post.Title = r.PostFormValue("title")
			post.Content = r.PostFormValue("postContent")
			post.CreationDate = time.Now().Format("January 2 15:04")
			m.db.InsertPostInfo(post)
			http.Redirect(w, r, "/main", http.StatusSeeOther)

		} else if r.Method == "GET" {
			tpl.ExecuteTemplate(w, "postCreate.html", nil)
		}

	}
}

// GetAllPosts ...
func (m *Multiplexer) GetAllPosts(w http.ResponseWriter) []*model.Post {

	posts, err := m.db.GetPosts()
	if err != nil {
		http.Error(w, "Something went wrong (Test Post)", http.StatusInternalServerError)
		return nil
	}
	return posts

}
