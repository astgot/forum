package database

import (
	"fmt"

	"github.com/astgot/forum/internal/model"
)

// InsertPostInfo ...
func (d *Database) InsertPostInfo(p *model.Post) (int64, error) {
	if err := d.Open(); err != nil {
		return -1, err
	}

	stmnt, err := d.db.Prepare("INSERT INTO Posts (user_id, author, title, content, creationDate) VALUES (?, ?, ?, ?, ?)")
	defer stmnt.Close()
	if err != nil {
		return -1, err
	}
	res, err := stmnt.Exec(p.UserID, p.Author, p.Title, p.Content, p.CreationDate)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

// GetPosts ...
func (d *Database) GetPosts() ([]*model.Post, error) {
	var posts []*model.Post
	if err := d.Open(); err != nil {
		return nil, err
	}
	res, err := d.db.Query("SELECT * FROM Posts")
	if err != nil {
		return nil, err
	}
	for res.Next() {
		post := model.NewPost()
		if err := res.Scan(&post.PostID, &post.UserID, &post.Author, &post.Title, &post.Content, &post.CreationDate); err != nil {
			// return nil, err
			fmt.Println(err, "query test")
		}
		posts = append(posts, post)

	}
	return posts, nil
}

// GetPostByPID ...
func (d *Database) GetPostByPID(pid int64) *model.Post {
	post := model.NewPost()
	if err := d.db.QueryRow("SELECT author, title, content, creationDate FROM Posts WHERE post_id = ?", pid).Scan(
		&post.Author,
		&post.Title,
		&post.Content,
		&post.CreationDate,
	); err != nil {
		return nil
	}
	return post
}
