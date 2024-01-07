package server

import "gitTimeline/pkg/models"

// TimelinePostStorage - a storage interface for the timeline
type TimelinePostStorage interface {
	CreatePost(post *models.PostContent) (postId string, err error)
	GetPost(id string) (string, error)
	GetPostTimeline(id string) ([]*models.PostVersionMetadata, error)
	GetPostVersion(id string, versionIdentifier string) (*models.PostContent, error)
	EditPost(postId string, postContent *models.PostContent) (string, error)
}
