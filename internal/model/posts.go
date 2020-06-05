package model

// Post ...
type Post struct {
	PostID       int64 `db:"postID"`
	UserID       int64 `db:"userID"`
	Title        string
	Author       string
	Content      string
	CreationDate string `db:"creationDate"`
}

// Category ...
type Category struct {
	CatID  int64
	PostID int64
	Name   string
}

// PostRating ...
type PostRating struct {
	PRID    int64 // Post Rating ID
	PostID  int64
	UID     int64 // userID
	Like    int64
	Dislike int64
}

// NewPost ...
func NewPost() *Post {
	return &Post{}
}

// NewCategory ...
func NewCategory() *Category {
	return &Category{}
}

// NewPostRating ...
func NewPostRating() *PostRating {
	return &PostRating{}
}
