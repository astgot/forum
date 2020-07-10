package controller

import (
	"fmt"
	"net/http"
	"strconv"
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
		// Check user authorization
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
		// Gathering post data
		post := model.NewPost()
		thread := model.NewThread()
		if r.Method == "POST" {
			r.ParseForm()
			post.UserID = u.ID
			post.Author = u.Firstname + " " + u.Lastname + " aka " + "\"" + u.Username + "\""
			post.Title = r.PostFormValue("title")
			post.Content = r.PostFormValue("postContent")
			thread.Name = r.PostFormValue("thread")
			threads := CheckNumberOfThreads(thread.Name)
			post.CreationDate = time.Now().Format("January 2 15:04")
			post.PostID, _ = m.db.InsertPostInfo(post)
			if post.PostID == -1 {
				fmt.Println("post.PostID == -1")
				http.Error(w, "Invalid input", http.StatusBadRequest)
				return
			}
			// If post has several threads, to this post will attach this info
			for _, threadName := range threads {
				m.db.InsertThreadInfo(threadName, post.PostID)
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)

		}
		// else if r.Method == "GET" {
		// 	tpl.ExecuteTemplate(w, "postCreate.html", nil)
		// }

	}
}

// PostView ... (/post?id=)
//(single post viewing -> to see comments, rate count OR if user is authenticated, he able to add comments and rate post here)
func (m *Multiplexer) PostView() http.HandlerFunc {

	type PostAttr struct {
		Threads  []*model.Thread
		Comments []*model.Comments
		//Likes, Dislikes
	}
	var singlePost struct {
		PostInfo []*PostAttr
		AuthUser *model.Users
		Post     *model.Post
	}
	return func(w http.ResponseWriter, r *http.Request) {

		id, errID := strconv.Atoi(r.URL.Query().Get("id"))
		// handling double post with id=1
		if id == 0 {
			id = -1
		}
		if errID != nil {
			fmt.Println("errID != nil")
			comment := r.URL.Query().Get("comment")
			if comment == "" && errID != nil {
				http.Error(w, "Invalid input", http.StatusBadRequest)
				return
			}
		}
		cookie, err := r.Cookie("authenticated")
		if err != nil {
			// If user is guest
			postAttr := &PostAttr{}
			singlePost.Post, err = m.db.GetPostByPID(int64(id))
			if err != nil {
				fmt.Println("Error on PostView() function")
				http.Error(w, "The post not found", http.StatusNotFound)
				return
			}
			postAttr.Comments, err = m.db.GetCommentsOfPost(int64(id))
			postAttr.Threads, _ = m.db.GetThreadOfPost(int64(id))
			singlePost.PostInfo = append(singlePost.PostInfo, postAttr)
			tpl.ExecuteTemplate(w, "postView.html", singlePost)
			return
		}
		postAttr := &PostAttr{}
		user, _ := m.db.GetUserByCookie(cookie.Value)
		singlePost.AuthUser = user
		singlePost.Post, err = m.db.GetPostByPID(int64(id))
		if err != nil {
			fmt.Println("Error on PostView() function")
			http.Error(w, "The post not found", http.StatusNotFound)
			return
		}
		postAttr.Comments, err = m.db.GetCommentsOfPost(int64(id))
		postAttr.Threads, _ = m.db.GetThreadOfPost(int64(id))
		singlePost.PostInfo = append(singlePost.PostInfo, postAttr)
		tpl.ExecuteTemplate(w, "postView.html", singlePost)
		if r.Method == "POST" {
			comment := model.NewComment()
			r.ParseForm()
			comment.Content = r.PostFormValue("comment")
			comment.CreationDate = time.Now().Format("Jan 2 15:04")
			comment.PostID = int64(id)
			comment.Author = user.Username // replace to comment's author?
			if ok := m.db.AddComment(comment); !ok {
				fmt.Println("AddComment error")
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
			// tpl.ExecuteTemplate(w, "postView.html", singlePost)
			// http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		// Add function to add Comments, rate Comments or Post

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
