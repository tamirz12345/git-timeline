package timeline

import (
	"fmt"
	"gitTimeline/pkg/models"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// GitStorage - a git storage implementation for the timeline
type GitStorage struct {
	logger                  *zerolog.Logger
	gitRepositoryCommunitor GitRepositoryCommunitor
}

// NewGitStorage - creates a new git storage
func NewGitStorage(gitRepositoryCommunitor GitRepositoryCommunitor, logger *zerolog.Logger) *GitStorage {
	return &GitStorage{
		logger:                  logger,
		gitRepositoryCommunitor: gitRepositoryCommunitor,
	}
}

// CreatePost - creates a new post in the repository, if nothing changed it will return an error
func (g *GitStorage) CreatePost(post *models.PostContent) (string, error) {
	// Generate a new UUID
	newUUID := uuid.New()

	// Convert the UUID to a string
	uuidStr := newUUID.String()
	filename := fmt.Sprintf("%s.json", uuidStr)
	// write to file
	postData := []byte(post.Body)
	commitId, err := g.gitRepositoryCommunitor.CommitFile(filename, post.Title, post.Username, postData)
	if err != nil {
		return "", fmt.Errorf("failed to write to file: %w", err)
	}
	if commitId == "" {
		return "", fmt.Errorf("nothing changed")
	}
	return uuidStr, nil
}

// GetPost - gets the latest version of a post
func (g *GitStorage) GetPost(id string) (string, error) {
	postPath := fmt.Sprintf("%s.json", id)
	content, _, err := g.gitRepositoryCommunitor.GetFileContent(postPath, "HEAD")
	if err != nil {
		g.logger.Error().Err(err).Str("postId", id).Msg("failed to get post content")
		return "", fmt.Errorf("failed to get target: %w", err)
	}
	return string(content), nil
}

// GetPostTimeline - gets the timeline of a post, meaning all the versions of the post leveraging git loog.
func (g *GitStorage) GetPostTimeline(id string) ([]*models.PostVersionMetadata, error) {
	postPath := fmt.Sprintf("%s.json", id)
	commits, err := g.gitRepositoryCommunitor.GetFileTimeline(postPath)
	if err != nil {
		g.logger.Error().Err(err).Str("postId", id).Msg("failed to get post timeline")
		return nil, fmt.Errorf("failed to get target: %w", err)
	}
	var postVersions []*models.PostVersionMetadata
	for _, commit := range commits.Commits {
		postVersions = append(postVersions, &models.PostVersionMetadata{
			VersionId: commit.CommitId,
			Title:     commit.Message,
			Date:      commit.Date,
			User:      commit.Author,
		})
	}
	return postVersions, nil
}

// GetPostVersion - gets a specific version of a post. if the version doesn't exist it will return an error. versionIdentifier is the commit hash
func (g *GitStorage) GetPostVersion(id string, versionIdentifier string) (*models.PostContent, error) {
	postPath := fmt.Sprintf("%s.json", id)
	content, commit, err := g.gitRepositoryCommunitor.GetFileContent(postPath, versionIdentifier)
	if err != nil {
		g.logger.Error().Err(err).Str("postId", id).Str("commitId", versionIdentifier).Msg("failed to get post content from past")
		return nil, fmt.Errorf("failed to get target: %w", err)
	}

	return &models.PostContent{
		Title:    commit.Message,
		Body:     string(content),
		Username: commit.Author,
	}, nil
}

// EditPost - edits a post in the repository, if nothing changed it will return an error
func (g *GitStorage) EditPost(postId string, post *models.PostContent) (string, error) {
	filename := fmt.Sprintf("%s.json", postId)
	postData := []byte(post.Body)
	commitId, err := g.gitRepositoryCommunitor.CommitFile(filename, post.Title, post.Username, postData)
	if err != nil {
		return "", fmt.Errorf("failed to write to file: %w", err)
	}
	if commitId == "" {
		g.logger.Info().Str("postId", postId).Msg("nothing changed")
		return "", fmt.Errorf("nothing changed")
	}
	return commitId, nil
}
