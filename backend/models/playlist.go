// backend/models/playlist.go

package models

import (
	"gorm.io/gorm"
)

// Playlist model
type Playlist struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Videos      []Video `gorm:"foreignKey:PlaylistID"`
}

// Video model
type Video struct {
	gorm.Model
	Title      string `json:"title"`
	URL        string `json:"url"`
	PlaylistID uint
}
