// backend/models/playlist.go

package models

import "gorm.io/gorm"

// Playlist model
type Playlist struct {
	gorm.Model
	Title            string            `json:"title"`
	Description      string            `json:"description"`
	Videos           []Video           `gorm:"foreignKey:PlaylistID"`
	ChannelID        uint              `json:"-"`
	Channel          Channel           `json:"channel"`
	IsFeatured       bool              `json:"isFeatured" gorm:"default:false"`
	IsPublic         bool              `json:"isPublic" gorm:"default:true"`
	FeaturedArtwork  string            `json:"featuredArtwork"`
	Tags             []string          `json:"tags" gorm:"type:varchar(255)[]"`
	IsPlayable       bool              `json:"isPlayable" gorm:"default:true"`
	PlayCount        uint              `json:"playCount" gorm:"default:0"`
	LikeCount        uint              `json:"likeCount" gorm:"default:0"`
	DislikeCount     uint              `json:"dislikeCount" gorm:"default:0"`
	Comments         []Comment         `gorm:"foreignKey:PlaylistID"`
	ShareCount       uint              `json:"shareCount" gorm:"default:0"`
	Followers        []User            `gorm:"many2many:user_playlist_followers;"`
	Contributors     []User            `gorm:"many2many:user_playlist_contributors;"`
	RelatedPlaylists []RelatedPlaylist `json:"relatedPlaylists" gorm:"foreignKey:PlaylistID"`
	TotalDuration    int               `json:"totalDuration" gorm:"default:0"`
	LastModified     int               `json:"lastModified" gorm:"autoUpdateTime"`
	PrivacySetting   PrivacySetting    `json:"privacySetting" gorm:"embedded"`
	Location         Location          `json:"location" gorm:"embedded"`
}

// UserPlaylistFollowers model for many-to-many relationship between users and playlists
type UserPlaylistFollowers struct {
	gorm.Model
	UserID     uint
	PlaylistID uint
}

// UserPlaylistContributors model for many-to-many relationship between users and playlists
type UserPlaylistContributors struct {
	gorm.Model
	UserID     uint
	PlaylistID uint
}

// RelatedPlaylist model for representing related playlists
type RelatedPlaylist struct {
	gorm.Model
	PlaylistID        uint
	RelatedPlaylistID uint
}

// PrivacySetting model for defining playlist privacy settings
type PrivacySetting struct {
	IsCollaborative bool `json:"isCollaborative" gorm:"default:false"`
	AllowComments   bool `json:"allowComments" gorm:"default:true"`
	// Add more privacy settings as needed
}
