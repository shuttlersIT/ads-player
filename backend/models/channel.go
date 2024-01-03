// backend/models/channel.go

package models

import (
	"time"

	"gorm.io/gorm"
)

// Channel model
type Channel struct {
	gorm.Model
	ID                 uint               `gorm:"primaryKey" json:"id"`
	Name               string             `gorm:"not null" json:"name"`
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

// TableName sets the table name for the Playlist model.
func (Channel) TableName() string {
	return "playlists"
}

// SocialMediaLinks model
type SocialMediaLinks struct {
	Facebook  string `json:"facebook"`
	Twitter   string `json:"twitter"`
	Instagram string `json:"instagram"`
	YouTube   string `json:"youtube"`
	// Add more social media links as needed
}

// TableName sets the table name for the Playlist model.
func (SocialMediaLinks) TableName() string {
	return "playlists"
}

// ContactInformation model
type ContactInformation struct {
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	// Add more contact information fields as needed
}

// TableName sets the table name for the Playlist model.
func (ContactInformation) TableName() string {
	return "playlists"
}

// Category model
type Category struct {
	gorm.Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Channels    []Channel `gorm:"many2many:channel_categories;"`
}

// TableName sets the table name for the Playlist model.
func (Category) TableName() string {
	return "playlists"
}

// ChannelDBModel provides methods for database operations related to Channel.
type ChannelDBModel struct {
	DB *gorm.DB
}

// NewChannelDBModel creates a new instance of ChannelDBModel.
func NewChannelDBModel(db *gorm.DB) *ChannelDBModel {
	return &ChannelDBModel{DB: db}
}

// CreateChannel creates a new channel in the database.
func (m *ChannelDBModel) CreateChannel(channel *Channel) error {
	return m.DB.Create(channel).Error
}

// GetChannelByID retrieves a channel by its ID from the database.
func (m *ChannelDBModel) GetChannelByID(channelID uint) (*Channel, error) {
	var channel Channel
	err := m.DB.First(&channel, channelID).Error
	return &channel, err
}

// GetChannels retrieves all channels from the database.
func (m *ChannelDBModel) GetChannels() ([]Channel, error) {
	var channels []Channel
	err := m.DB.Find(&channels).Error
	return channels, err
}

// UpdateChannel updates the information of a channel in the database.
func (m *ChannelDBModel) UpdateChannel(channelID uint, updatedChannel *Channel) error {
	var channel Channel
	err := m.DB.First(&channel, channelID).Error
	if err != nil {
		return err
	}

	channel.Name = updatedChannel.Name

	return m.DB.Save(&channel).Error
}

// DeleteChannel deletes a channel from the database.
func (m *ChannelDBModel) DeleteChannel(channelID uint) error {
	return m.DB.Delete(&Channel{}, channelID).Error
}
