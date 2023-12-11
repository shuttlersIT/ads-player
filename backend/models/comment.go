// backend/models/comment.go

package models

import "gorm.io/gorm"

// Comment model
type Comment struct {
	gorm.Model
	UserID           uint              `json:"-"`
	User             User              `json:"user"`
	VideoID          uint              `json:"-"`
	Video            Video             `json:"video"`
	Text             string            `json:"text"`
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

// RelatedComment model for representing related comments
type RelatedComment struct {
	gorm.Model
	CommentID        uint
	RelatedCommentID uint
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

// Hashtag model for representing hashtags in comments
type Hashtag struct {
	gorm.Model
	Name     string    `json:"name"`
	Comments []Comment `gorm:"many2many:comment_hashtags;"`
	// Add more hashtag-related fields as needed
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

// Emoji model for representing emojis in comments
type Emoji struct {
	gorm.Model
	Name     string    `json:"name"`
	Comments []Comment `gorm:"many2many:comment_emojis;"`
	// Add more emoji-related fields as needed
}
