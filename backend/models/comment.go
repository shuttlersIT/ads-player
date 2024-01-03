// backend/models/comment.go

package models

import (
	"time"

	"gorm.io/gorm"
)

// Comment model
type Comment struct {
	gorm.Model
	ID               uint              `gorm:"primaryKey" json:"id"`
	UserID           uint              `gorm:"not null" json:"user_id"`
	User             User              `gorm:"foreignKey:UserID" json:"user"`
	VideoID          uint              `gorm:"not null" json:"video_id"`
	Video            Video             `json:"video"`
	Content          string            `gorm:"not null" json:"content"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	DeletedAt        time.Time         `gorm:"index" json:"-"`
	Text             string            `gorm:"not null" json:"text"`
	Likes            uint              `json:"likes" gorm:"default:0"`
	Dislikes         uint              `json:"dislikes" gorm:"default:0"`
	Replies          []Comment         `gorm:"foreignKey:ParentCommentID"`
	ParentCommentID  uint              `json:"-"`
	IsReported       bool              `json:"isReported" gorm:"default:false"`
	ReportsCount     uint              `json:"reportsCount" gorm:"default:0"`
	IsEdited         bool              `json:"isEdited" gorm:"default:false"`
	EditedTimestamp  int               `json:"editedTimestamp"`
	IsDeleted        bool              `json:"isDeleted" gorm:"default:false"`
	DeletedTimestamp int               `json:"deletedTimestamp"`
	MediaAttachments []MediaAttachment `json:"mediaAttachments"`
	UserMentions     []User            `gorm:"many2many:comment_user_mentions;"`
	Hashtags         []Hashtag         `gorm:"many2many:comment_hashtags;"`
	Location         Location          `json:"location" gorm:"embedded"`
	Emojis           []Emoji           `gorm:"many2many:comment_emojis;"`
}

// TableName sets the table name for the Playlist model.
func (Comment) TableName() string {
	return "playlists"
}

// RelatedComment model for representing related comments
type RelatedComment struct {
	gorm.Model
	CommentID        uint
	RelatedCommentID uint
}

// TableName sets the table name for the Playlist model.
func (RelatedComment) TableName() string {
	return "playlists"
}

// MediaAttachment model for representing attached media in a comment
type MediaAttachment struct {
	gorm.Model
	CommentID uint `json:"-"`
	Comment   Comment
	URL       string `json:"url"`
	Type      string `json:"type"`
	Caption   string `json:"caption"`
	// Add more media attachment fields as needed
}

// TableName sets the table name for the Playlist model.
func (MediaAttachment) TableName() string {
	return "playlists"
}

// Hashtag model for representing hashtags in comments
type CommentsHashtag struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey" json:"id"`
	Tag       string    `json:"tag"`
	Tags      string    `json:"tags"` // Added Tags field
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Comments  []Comment `gorm:"many2many:comment_hashtags;"`
	// Add more hashtag-related fields as needed
}

// TableName sets the table name for the Playlist model.
func (CommentsHashtag) TableName() string {
	return "playlists"
}

// Location model for representing geographic location in comments
type Location struct {
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	Name             string  `json:"name"`
	Address          string  `json:"address"`
	City             string  `json:"city"`
	State            string  `json:"state"`
	Country          string  `json:"country"`
	ZipCode          string  `json:"zipCode"`
	Region           string  `json:"region"`
	PlaceID          string  `json:"placeID"`
	FormattedAddress string  `json:"formattedAddress"`
	// Add more location-related fields as needed
}

// TableName sets the table name for the Playlist model.
func (Location) TableName() string {
	return "playlists"
}

// Emoji model for representing emojis in comments
type Emoji struct {
	gorm.Model
	Name     string    `json:"name"`
	Comments []Comment `gorm:"many2many:comment_emojis;"`
	// Add more emoji-related fields as needed
}

// TableName sets the table name for the Playlist model.
func (Emoji) TableName() string {
	return "playlists"
}

// CommentDBModel provides methods for database operations related to Comment.
type CommentDBModel struct {
	DB *gorm.DB
}

// NewCommentDBModel creates a new instance of CommentDBModel.
func NewCommentDBModel(db *gorm.DB) *CommentDBModel {
	return &CommentDBModel{DB: db}
}

// GetCommentsByVideoID retrieves comments for a specific video from the database.
func (m *CommentDBModel) GetCommentsByVideoID(videoID uint) ([]Comment, error) {
	var comments []Comment
	err := m.DB.Where("video_id = ?", videoID).Find(&comments).Error
	return comments, err
}

// CreateComment creates a new comment in the database.
func (m *CommentDBModel) CreateComment(comment *Comment) error {
	return m.DB.Create(comment).Error
}

// UpdateComment updates the information of a comment in the database.
func (m *CommentDBModel) UpdateComment(commentID uint, updatedComment *Comment) error {
	var comment Comment
	err := m.DB.First(&comment, commentID).Error
	if err != nil {
		return err
	}

	comment.Content = updatedComment.Content

	return m.DB.Save(&comment).Error
}

// DeleteComment deletes a comment from the database.
func (m *CommentDBModel) DeleteComment(commentID uint) error {
	return m.DB.Delete(&Comment{}, commentID).Error
}
