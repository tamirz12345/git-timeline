package models

import "time"

type PostContent struct {
	Body     string `json:"body"`
	Title    string `json:"title"`
	Username string `json:"username"`
}

type PostMetadata struct {
	PostID         string `json:"postId"`
	Title          string `json:"title"`
	Date           string `json:"date"`
	VersionsNumber int    `json:"versionsNumber"`
	Username       string `json:"username"`
}

type PostVersionMetadata struct {
	VersionId string    `json:"versionId"`
	Title     string    `json:"title"`
	Date      time.Time `json:"date"`
	User      string    `json:"user"`
}
