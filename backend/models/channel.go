// backend/models/channel.go

package models

import (
	"time"

	"gorm.io/gorm"
)

// Channel model
type Channel struct {
	gorm.Model
	Name               string             `json:"name"`
	Description        string             `json:"description"`
	OwnerID            uint               `json:"-"`
	Owner              User               `json:"owner"`
	ProfilePicture     string             `json:"profilePicture"`
	BannerImage        string             `json:"bannerImage"`
	Subscribers        []User             `gorm:"many2many:user_subscriptions;"`
	Playlists          []Playlist         `json:"playlists"`
	FeaturedVideos     []Video            `json:"featuredVideos" gorm:"many2many:channel_featured_videos;"`
	Uploads            []Video            `json:"uploads"`
	Categories         []Category         `json:"categories" gorm:"many2many:channel_categories;"`
	SocialMediaLinks   SocialMediaLinks   `json:"socialMediaLinks" gorm:"embedded"`
	ContactInformation ContactInformation `json:"contactInformation" gorm:"embedded"`
	Website            string             `json:"website"`
	About              string             `json:"about"`
	SubscribersCount   uint               `json:"subscribersCount" gorm:"default:0"`
	ViewsCount         uint               `json:"viewsCount" gorm:"default:0"`
	IsVerified         bool               `json:"isVerified" gorm:"default:false"`
	JoinDate           time.Time          `json:"joinDate"`
	LastUploadDate     time.Time          `json:"lastUploadDate"`
	MonetarySupportURL string             `json:"monetarySupportURL"`
}

// SocialMediaLinks model
type SocialMediaLinks struct {
	Facebook  string `json:"facebook"`
	Twitter   string `json:"twitter"`
	Instagram string `json:"instagram"`
	YouTube   string `json:"youtube"`
	// Add more social media links as needed
}

// ContactInformation model
type ContactInformation struct {
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	// Add more contact information fields as needed
}

// Category model
type Category struct {
	gorm.Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Channels    []Channel `gorm:"many2many:channel_categories;"`
}
