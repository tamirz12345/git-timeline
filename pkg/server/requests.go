package server

import (
	"gitTimeline/pkg/models"
)

// CreatePostRequest - a request for create post
type CreatePostRequest models.PostContent

// UpdatePostRequest - a request for update post
type UpdatePostRequest struct {
	Body  string `json:"body"`
	Title string `json:"title"`
	User  string `json:"username"`
}
