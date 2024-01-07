package server

import "gitTimeline/pkg/models"

// PostStorage - a storage interface for the metadata
type PostMetadataStorage interface {
	CreatePostMetadata(post *models.PostMetadata) error
	UpdatePostMetadata(post *models.PostMetadata) error
	GetPostMetadata(postId string) (*models.PostMetadata, error)
	GetAllPostsMetadata() ([]*models.PostMetadata, error)
}
