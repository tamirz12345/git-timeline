package server

// CreatePostResponse - a response for create post
type CreatePostResponse struct {
	PostId string `json:"postId"`
}

// CreateVersionResponse - a response for create version
type CreateVersionResponse struct {
	VersionId string `json:"versionId"`
}

// GetPostsResponse - a response for get posts request
type GetPostResponse struct {
	Body           string `json:"body"`
	Title          string `json:"title"`
	VersionsNumber int    `json:"versionsNumber"`
	Username       string `json:"username"`
	Date           string `json:"date"`
}
