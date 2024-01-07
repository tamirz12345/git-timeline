package metadata

import (
	"fmt"
	"gitTimeline/pkg/models"

	"github.com/rs/zerolog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlStorage struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewMysqlStorage - creates a new mysql storage, if the table doesn't exist it will create it
func NewMysqlStorage(dbHost, dbUser, dbPass, dbName string, logger *zerolog.Logger) (*MysqlStorage, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	db.AutoMigrate(&models.PostMetadata{})

	logger.Info().Msg("ensuring posts table exists")
	return &MysqlStorage{db: db, logger: logger}, nil
}

// CreatePostMetadata - add new post to the db
func (m *MysqlStorage) CreatePostMetadata(post *models.PostMetadata) error {
	return m.db.Create(post).Error
}

// UpdatePostMetadata - updates a post metadata
func (m *MysqlStorage) UpdatePostMetadata(post *models.PostMetadata) error {
	return m.db.Model(&models.PostMetadata{}).Where("post_id = ?", post.PostID).Updates(post).Error
}

// GetPostMetadata - gets a post metadata by id
func (m *MysqlStorage) GetPostMetadata(postId string) (*models.PostMetadata, error) {
	var post models.PostMetadata
	if err := m.db.Where("post_id = ?", postId).First(&post).Error; err != nil {
		m.logger.Error().Err(err).Str("postId", postId).Msg("failed to get post metadata")
		return nil, fmt.Errorf("failed to get post metadata: %w", err)
	}
	return &post, nil
}

// GetAllPostsMetadata - gets all posts metadata
func (m *MysqlStorage) GetAllPostsMetadata() ([]*models.PostMetadata, error) {
	var posts []*models.PostMetadata
	if err := m.db.Find(&posts).Error; err != nil {
		m.logger.Error().Err(err).Msg("failed to get all posts metadata")
		return nil, fmt.Errorf("failed to get all posts metadata: %w", err)
	}
	return posts, nil
}
