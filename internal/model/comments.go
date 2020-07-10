package model

// Comments ...
type Comments struct {
	CommentID    int64  `db:"ID"`
	PostID       int64  `db:"postID"`
	Author       string `db:"author"`
	Content      string `db:"content"`
	CreationDate string `db:"creationDate"`
}

// NewComment ...
func NewComment() *Comments {
	return &Comments{}
}
