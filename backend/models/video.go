// backend/models/video.go

package models

import "gorm.io/gorm"

// Video model
type Video struct {
	gorm.Model
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	URL             string         `json:"url"`
	ThumbnailURL    string         `json:"thumbnailUrl"`
	Views           uint           `json:"views" gorm:"default:0"`
	Likes           uint           `json:"likes" gorm:"default:0"`
	Dislikes        uint           `json:"dislikes" gorm:"default:0"`
	Comments        []Comment      `gorm:"foreignKey:VideoID"`
	PlaylistID      uint           `json:"-"`
	Playlist        Playlist       `json:"playlist"`
	IsAdvertisement bool           `json:"isAdvertisement" gorm:"default:false"`
	Duration        int            `json:"duration"` // Duration in seconds
	Order           int            `json:"order" gorm:"default:0"`
	Tags            []string       `json:"tags" gorm:"type:varchar(255)[]"`
	UploadDate      int            `json:"uploadDate" gorm:"autoCreateTime"`
	Uploader        User           `json:"uploader"`
	Category        Category       `json:"category"`
	PrivacySetting  PrivacySetting `json:"privacySetting" gorm:"embedded"`
	CommentsEnabled bool           `json:"commentsEnabled" gorm:"default:true"`
	RelatedVideos   []RelatedVideo `json:"relatedVideos" gorm:"foreignKey:VideoID"`
	// Add more video-related fields as needed
}

// RelatedVideo model for representing related videos
type RelatedVideo struct {
	gorm.Model
	VideoID        uint
	RelatedVideoID uint
}
