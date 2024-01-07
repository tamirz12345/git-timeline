package server

import (
	"net/http"
	"strconv"
	"time"

	"gitTimeline/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// Server - a server for the timeline
type Server struct {
	engine          *gin.Engine
	logger          *zerolog.Logger
	port            int
	timelineStorage TimelinePostStorage
	metadataStorage PostMetadataStorage
}

// GetPostsResponse - a response for get posts request
func NewServer(logger *zerolog.Logger, postStorage TimelinePostStorage, metadataStorage PostMetadataStorage, port int) *Server {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	server := &Server{
		engine:          engine,
		logger:          logger,
		timelineStorage: postStorage,
		port:            port,
		metadataStorage: metadataStorage,
	}

	server.registerRoutes()

	return server
}

// registerRoutes - registers the routes for the server
func (s *Server) registerRoutes() {
	r := s.engine
	r.GET("/api/posts", s.HandleGetPosts)
	r.GET("/api/post/:id", s.HandleGetPost)
	r.GET("/api/post/:id/timeline", s.HandleGetTimeline)
	r.GET("/api/post/:id/version/:versionIdentifier", s.HandleGetVersion)
	r.POST("/api/post", s.HandleCreatePost)
	r.PUT("/api/post/:id", s.HandleEditPost)
}

// Start - starts the server
func (s *Server) Start() error {
	s.logger.Info().Int("port", s.port).Msg("running server")
	return s.engine.Run(":" + strconv.Itoa(s.port))
}

// HandleGetPosts - handles get posts request.
func (s *Server) HandleGetPosts(c *gin.Context) {
	posts, err := s.metadataStorage.GetAllPostsMetadata()
	if err != nil {
		s.logger.Error().Err(err).Msg("error getting posts")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// HandleGetPost - handles get post request. Gets the latest version of the post. In case the post doesn't exist returns 404.
func (s *Server) HandleGetPost(c *gin.Context) {
	postID := c.Param("id")
	postMetadata, err := s.metadataStorage.GetPostMetadata(postID)
	if err != nil || postMetadata == nil {
		s.logger.Error().Err(err).Msg("post not exist")
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	postContent, err := s.timelineStorage.GetPost(postID)
	if err != nil {
		s.logger.Error().Err(err).Msg("error getting post")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch post content"})
		return
	}
	c.JSON(http.StatusOK, GetPostResponse{
		Body:           postContent,
		Title:          postMetadata.Title,
		VersionsNumber: postMetadata.VersionsNumber,
		Username:       postMetadata.Username,
		Date:           postMetadata.Date,
	})
}

// HandleGetTimeline - handles get timeline request. Gets all the versions of the post. In case the post doesn't exist returns 404.
// For each version it returns the version identifier, version title, the date of the version and the username of the user who created the version.
func (s *Server) HandleGetTimeline(c *gin.Context) {
	postId := c.Param("id")
	postMetadata, err := s.metadataStorage.GetPostMetadata(postId)
	if err != nil || postMetadata == nil {
		s.logger.Error().Err(err).Msg("post not exist")
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	postTimeline, err := s.timelineStorage.GetPostTimeline(postId)
	if err != nil {
		s.logger.Error().Err(err).Msg("error getting timeline")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch post timeline"})
		return
	}
	c.JSON(http.StatusOK, postTimeline)
}

// HandleGetVersion - handles get version request. Gets a specific version of the post.
// In case the post doesn't or the version doesn't exist returns 404.
func (s *Server) HandleGetVersion(c *gin.Context) {
	postId := c.Param("id")
	versionId := c.Param("versionIdentifier")
	if versionId == "" || postId == "" {
		s.logger.Warn().Str("postId", postId).Str("versionId", versionId).Msg("invalid request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	postContent, err := s.timelineStorage.GetPostVersion(postId, versionId)
	if err != nil {
		s.logger.Error().Err(err).Msg("error getting version")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch post version"})
		return
	}
	c.JSON(http.StatusOK, postContent)
}

// HandleCreatePost - handles create post request. Creates a new post and returns the post id.
func (s *Server) HandleCreatePost(c *gin.Context) {
	var createPostRequest CreatePostRequest
	if err := c.ShouldBindJSON(&createPostRequest); err != nil {
		s.logger.Error().Err(err).Msg("error decoding request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Create a new post using the timeline storage
	postId, err := s.timelineStorage.CreatePost(&models.PostContent{
		Body:     createPostRequest.Body,
		Title:    createPostRequest.Title,
		Username: createPostRequest.Username,
	})
	if err != nil {
		s.logger.Error().Err(err).Msg("error creating post")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	// Create post metadata using the metadata storage
	err = s.metadataStorage.CreatePostMetadata(&models.PostMetadata{
		PostID:         postId,
		Title:          createPostRequest.Title,
		Username:       createPostRequest.Username,
		VersionsNumber: 1,
		Date:           time.Now().Format(time.RFC3339),
	})
	if err != nil {
		s.logger.Error().Err(err).Msg("error creating post metadata")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post metadata"})
		return
	}

	c.JSON(http.StatusOK, CreatePostResponse{
		PostId: postId,
	})
}

// HandleEditPost - handles edit post request. Edits an existing post and returns the version id.
func (s *Server) HandleEditPost(c *gin.Context) {
	var updatePostRequest UpdatePostRequest
	if err := c.ShouldBindJSON(&updatePostRequest); err != nil {
		s.logger.Error().Err(err).Msg("error decoding request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	postID := c.Param("id")

	postMetadata, err := s.metadataStorage.GetPostMetadata(postID)
	if err != nil || postMetadata == nil {
		s.logger.Error().Err(err).Msg("error getting post metadata")
		c.JSON(http.StatusNotFound, gin.H{"error": "Post metadata not found"})
		return
	}

	versionId, err := s.timelineStorage.EditPost(postID, &models.PostContent{
		Body:     updatePostRequest.Body,
		Title:    updatePostRequest.Title,
		Username: updatePostRequest.User,
	})

	if err != nil {
		if err.Error() == "nothing changed" {
			s.logger.Info().Str("postId", postID).Msg("nothing changed")
			c.JSON(http.StatusNotModified, gin.H{"message": "No changes detected"})
			return
		}
		s.logger.Error().Err(err).Msg("error editing post")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit post"})
		return
	}

	err = s.metadataStorage.UpdatePostMetadata(&models.PostMetadata{
		PostID:         postID,
		Title:          updatePostRequest.Title,
		Username:       updatePostRequest.User,
		VersionsNumber: postMetadata.VersionsNumber + 1,
		Date:           time.Now().Format(time.RFC3339),
	})

	if err != nil {
		s.logger.Error().Err(err).Msg("error updating post metadata")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post metadata"})
		return
	}

	c.JSON(http.StatusOK, CreateVersionResponse{
		VersionId: versionId,
	})
}
