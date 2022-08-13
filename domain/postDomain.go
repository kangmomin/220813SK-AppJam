package domain

import "time"

type Post struct {
	PostId      int       `json:"post_id"`
	OwnerId     int       `json:"owner_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	OwnerName   string    `json:"owner_name"`
	Created     time.Time `json:"created"`
}

type PostList struct {
	PostId   int    `json:"post_id"`
	Title    string `json:"titile"`
	UserName string `json:"user_name"`
	Created  int    `json:"created"`
}

type WritePost struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
