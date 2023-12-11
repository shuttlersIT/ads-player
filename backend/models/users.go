// backend/models/user.go

package models

import (
	"gorm.io/gorm"
)

// User model
type User struct {
	gorm.Model
	Username       string     `json:"username"`
	Email          string     `json:"email" gorm:"unique"`
	Password       string     `json:"-"`
	ProfilePicture string     `json:"profilePicture"`
	FirstName      string     `json:"firstName"`
	LastName       string     `json:"lastName"`
	Bio            string     `json:"bio"`
	Subscriptions  []Channel  `gorm:"many2many:user_subscriptions;"`
	Playlists      []Playlist `json:"playlists"`
	WatchHistory   []Video    `json:"watchHistory" gorm:"many2many:user_watch_history;"`
	LikedVideos    []Video    `json:"likedVideos" gorm:"many2many:user_liked_videos;"`
	DislikedVideos []Video    `json:"dislikedVideos" gorm:"many2many:user_disliked_videos;"`
}
