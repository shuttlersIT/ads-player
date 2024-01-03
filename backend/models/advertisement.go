// backend/models/advertisement.go

package models

import (
	"time"

	"gorm.io/gorm"
)

// Advertisement model
type Advertisement struct {
	gorm.Model
	ID                    uint                   `json:"id" gorm:"primaryKey"`
	Name                  string                 `json:"name"`
	PlaylistID            uint                   `json:"-"`
	Playlist              Playlist               `json:"playlist"`
	ContentURL            string                 `json:"content_url"`
	Title                 string                 `json:"title"`
	Description           string                 `json:"description"`
	Duration              uint                   `json:"duration"` // Duration in seconds
	ScheduledAt           time.Time              `json:"scheduledAt"`
	Played                bool                   `json:"played" gorm:"default:false"`
	ClickThroughURL       string                 `json:"clickThroughURL"`
	Analytics             AdvertisementAnalytics `json:"analytics" gorm:"embedded"`
	IsFeatured            bool                   `json:"isFeatured" gorm:"default:false"`
	IsPublic              bool                   `json:"isPublic" gorm:"default:true"`
	Tags                  []string               `json:"tags" gorm:"type:varchar(255)[]"`
	LikeCount             uint                   `json:"likeCount" gorm:"default:0"`
	DislikeCount          uint                   `json:"dislikeCount" gorm:"default:0"`
	Comments              []Comment              `gorm:"foreignKey:AdvertisementID"`
	ShareCount            uint                   `json:"shareCount" gorm:"default:0"`
	Followers             []User                 `gorm:"many2many:user_advertisement_followers;"`
	Contributors          []User                 `gorm:"many2many:user_advertisement_contributors;"`
	RelatedAds            []RelatedAd            `json:"relatedAds" gorm:"foreignKey:AdvertisementID"`
	LastModified          int                    `json:"lastModified" gorm:"autoUpdateTime"`
	PrivacySetting        PrivacySetting         `json:"privacySetting" gorm:"embedded"`
	Location              Location               `json:"location" gorm:"embedded"`
	VideoQuality          string                 `json:"videoQuality"`
	AudioQuality          string                 `json:"audioQuality"`
	Caption               string                 `json:"caption"`
	Language              string                 `json:"language"`
	Label                 string                 `json:"label"`
	CreatedAt             time.Time              `json:"createdAt"`
	DeletedAt             time.Time              `gorm:"index" json:"-"`
	UpdatedAt             time.Time              `json:"updatedAt"`
	TargetAudience        string                 `json:"targetAudience"`
	MatureContent         bool                   `json:"matureContent"`
	ThumbnailURL          string                 `json:"thumbnailURL"`
	ExternalLinks         []ExternalLink         `json:"externalLinks" gorm:"foreignKey:AdvertisementID"`
	MediaAttachments      []MediaAttachment      `json:"mediaAttachments" gorm:"foreignKey:AdvertisementID"`
	AdvertisementHashtags []AdvertisementHashtag `json:"hashtags" gorm:"foreignKey:AdvertisementID"`
	StartDate             time.Time              `json:"start_date"`
	EndDate               time.Time              `json:"end_date"`
	VideoID               uint                   `json:"video_id" gorm:"not null"`
	CreatorUserID         uint                   `json:"creator_user_id"`
	PlayCount             int                    `json:"playCount" gorm:"default:0"`
}

// TableName sets the table name for the Advertisement model.
func (Advertisement) TableName() string {
	return "advertisements"
}

// AdvertisementAnalytics struct for tracking advertisement analytics
type AdvertisementAnalytics struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	AdvertisementID      uint      `gorm:"index" json:"advertisement_id"`
	Views                uint      `json:"views" gorm:"default:0"`
	Clicks               uint      `json:"clicks" gorm:"default:0"`
	Comments             uint      `json:"comments" gorm:"default:0"`
	Shares               uint      `json:"shares" gorm:"default:0"`
	Likes                uint      `json:"likes" gorm:"default:0"`
	Dislikes             uint      `json:"dislikes" gorm:"default:0"`
	Followers            uint      `json:"followers" gorm:"default:0"`
	PlayCount            uint      `json:"playCount" gorm:"default:0"`
	TotalDurationWatched int       `json:"totalDurationWatched" gorm:"default:0"`
	ClickThroughCount    uint      `json:"clickThroughCount" gorm:"default:0"`
	ConversionRate       float64   `json:"conversionRate" gorm:"default:0.0"`
	Conversions          uint      `json:"conversions"` // Added Conversions field
	Revenue              float64   `json:"revenue"`     // Added Revenue field
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	// Add more analytics fields as needed
}

// TableName sets the table name for the Advertisement model.
func (AdvertisementAnalytics) TableName() string {
	return "advertisements"
}

// PrivacySetting struct for advertisement privacy settings
type AdvertisementPrivacySetting struct {
	IsPublic       bool   `json:"isPublic" gorm:"default:true"`
	Password       string `json:"password,omitempty"`
	AllowComments  bool   `json:"allowComments" gorm:"default:true"`
	AllowDownloads bool   `json:"allowDownloads" gorm:"default:true"`
	AllowEmbedding bool   `json:"allowEmbedding" gorm:"default:true"`
}

// TableName sets the table name for the Advertisement model.
func (AdvertisementPrivacySetting) TableName() string {
	return "advertisements"
}

// Location struct for advertisement location details
type AdvertisementLocation struct {
	Latitude      float64 `json:"latitude,omitempty"`
	Longitude     float64 `json:"longitude,omitempty"`
	Country       string  `json:"country,omitempty"`
	City          string  `json:"city,omitempty"`
	State         string  `json:"state,omitempty"`
	ZipCode       string  `json:"zipCode,omitempty"`
	StreetAddress string  `json:"streetAddress,omitempty"`
}

// TableName sets the table name for the Advertisement model.
func (AdvertisementLocation) TableName() string {
	return "advertisements"
}

// ExternalLink struct for storing external links related to the advertisement
type ExternalLink struct {
	gorm.Model
	AdvertisementID uint   `json:"-"`
	URL             string `json:"url"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	IsFeatured      bool   `json:"isFeatured" gorm:"default:false"`
	IsAffiliate     bool   `json:"isAffiliate" gorm:"default:false"`
	Label           string `json:"label"`
	// Add more fields as needed
}

// TableName sets the table name for the Advertisement model.
func (ExternalLink) TableName() string {
	return "advertisements"
}

// MediaAttachment struct for storing media attachments related to the advertisement
type AdvertisementMediaAttachment struct {
	gorm.Model
	AdvertisementID uint   `json:"-"`
	URL             string `json:"url"`
	Type            string `json:"type"`
	Caption         string `json:"caption"`
	AltText         string `json:"altText"`
	IsPrimary       bool   `json:"isPrimary" gorm:"default:false"`
	Order           int    `json:"order" gorm:"default:0"`
	// Add more fields as needed
}

// TableName sets the table name for the Advertisement model.
func (AdvertisementMediaAttachment) TableName() string {
	return "advertisements"
}

// Hashtag struct for storing hashtags related to the advertisement
type AdvertisementHashtag struct {
	gorm.Model
	ID              uint      `gorm:"primaryKey" json:"id"`
	Tag             string    `json:"tag"`
	Tags            string    `json:"tags"` // Added Tags field
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	AdvertisementID uint      `json:"-"`
	Text            string    `json:"text"`
	// Add more fields as needed
}

// TableName sets the table name for the Advertisement model.
func (AdvertisementHashtag) TableName() string {
	return "advertisements"
}

// Hashtag represents a hashtag entity
type Hashtag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Tag       string    `json:"tag"`
	Tags      string    `json:"tags"` // Added Tags field
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName sets the table name for the Advertisement model.
func (Hashtag) TableName() string {
	return "advertisements"
}

// RelatedAd struct for storing related advertisements
type RelatedAd struct {
	gorm.Model
	AdvertisementID        uint `json:"-"`
	RelatedAdvertisementID uint `json:"relatedAdvertisementID"`
	Order                  int  `json:"order" gorm:"default:0"`
}

// TableName sets the table name for the Advertisement model.
func (RelatedAd) TableName() string {
	return "advertisements"
}

// AdvertisementModel handles database operations for Advertisement
type AdvertisementDBModel struct {
	DB *gorm.DB
}

// NewAdvertisementModel creates a new instance of AdvertisementModel
func NewAdvertisementDBModel(db *gorm.DB) *AdvertisementDBModel {
	return &AdvertisementDBModel{
		DB: db,
	}
}

// AdvertisementModel is the model interface for advertisement-related operations
type AdvertisementModel interface {
	GetAllAdvertisements() ([]Advertisement, error)
	GetAdvertisementByID(id uint) (*Advertisement, error)
	CreateAdvertisement(ad *Advertisement) error
	UpdateAdvertisement(ad *Advertisement) error
	DeleteAdvertisement(id uint) error
	GetCurrentAdvertisement(id uint) (*Advertisement, error)
	GetNextAdvertisement(id uint) (*Advertisement, error)
	GetAdvertisementsByPlaylistID(playlistID uint) ([]Advertisement, error)
	GetNextAdvertisementForPlaylist(playlistID uint) (*Advertisement, error)
	MarkAdvertisementAsPlayed(id uint) error
	UpdateAdvertisementPlayStatus(id uint, played bool) error
	GetPlaybackRate() (float64, error)
	GetVideoIDForAdvertisement(id uint) (uint, error)
	IncrementPlayCount(advertisementID uint)
}

// GetAdvertisementByID retrieves an advertisement by its ID.
func (m *AdvertisementDBModel) GetAdvertisementByID(id uint) (*Advertisement, error) {
	var advertisement Advertisement
	err := m.DB.Where("id = ?", id).First(&advertisement).Error
	return &advertisement, err
}

// UpdateAdvertisement updates the details of an existing advertisement.
func (m *AdvertisementDBModel) UpdateAdvertisement(advertisement *Advertisement) error {
	return m.DB.Save(advertisement).Error
}

// DeleteAdvertisement deletes an advertisement from the database.
func (m *AdvertisementDBModel) DeleteAdvertisement(id uint) error {
	return m.DB.Delete(&Advertisement{}, id).Error
}

// GetAllAdvertisements retrieves all advertisements from the database.
func (m *AdvertisementDBModel) GetAllAdvertisements() ([]Advertisement, error) {
	var advertisements []Advertisement
	err := m.DB.Find(&advertisements).Error
	return advertisements, err
}

/*-----------------------------------------------------------------------------------------------------------------------*/

// GetNextAdvertisementForPlaylist fetches the next advertisement to play for a playlist
func (am *AdvertisementDBModel) GetNextAdvertisementForPlaylist(playlistID uint) (*Advertisement, error) {
	var advertisement Advertisement
	if err := am.DB.Where("playlist_id = ? AND scheduled_at <= ? AND played = ?", playlistID, time.Now(), false).Order("scheduled_at").First(&advertisement).Error; err != nil {
		return nil, err
	}

	return &advertisement, nil
}

// MarkAdvertisementAsPlayed marks an advertisement as played (update status or other relevant fields)
func (am *AdvertisementDBModel) MarkAdvertisementAsPlayed(advertisementID uint) error {
	var advertisement Advertisement
	if err := am.DB.First(&advertisement, advertisementID).Error; err != nil {
		return err
	}

	// Add your logic to update status or other relevant fields
	advertisement.Played = true

	if err := am.DB.Save(&advertisement).Error; err != nil {
		return err
	}

	return nil
}

// GetAdvertisementByID fetches an advertisement by its ID
func (am *AdvertisementDBModel) GetAdvertisementByID2(advertisementID uint) (*Advertisement, error) {
	var advertisement Advertisement
	if err := am.DB.Preload("Playlist").Preload("Comments").Preload("Followers").Preload("Contributors").Preload("RelatedAds").First(&advertisement, advertisementID).Error; err != nil {
		return nil, err
	}
	return &advertisement, nil
}

// GetAllAdvertisements fetches all advertisements
func (am *AdvertisementDBModel) GetAllAdvertisements2() ([]Advertisement, error) {
	var advertisements []Advertisement
	if err := am.DB.Preload("Playlist").Preload("Comments").Preload("Followers").Preload("Contributors").Preload("RelatedAds").Find(&advertisements).Error; err != nil {
		return nil, err
	}
	return advertisements, nil
}

// CreateAdvertisement creates a new advertisement
func (am *AdvertisementDBModel) CreateAdvertisement(advertisement *Advertisement) error {
	if err := am.DB.Create(advertisement).Error; err != nil {
		return err
	}
	return nil
}

// UpdateAdvertisement updates an existing advertisement
func (am *AdvertisementDBModel) UpdateAdvertisement2(advertisement *Advertisement) error {
	if err := am.DB.Save(advertisement).Error; err != nil {
		return err
	}
	return nil
}

// DeleteAdvertisement deletes an advertisement by its ID
func (am *AdvertisementDBModel) DeleteAdvertisement2(advertisementID uint) error {
	if err := am.DB.Delete(&Advertisement{}, advertisementID).Error; err != nil {
		return err
	}
	return nil
}

// GetAdvertisementsByPlaylistID fetches all advertisements for a specific playlist
func (am *AdvertisementDBModel) GetAdvertisementsByPlaylistID(playlistID uint) ([]Advertisement, error) {
	var advertisements []Advertisement
	if err := am.DB.Where("playlist_id = ?", playlistID).Preload("Playlist").Preload("Comments").Preload("Followers").Preload("Contributors").Preload("RelatedAds").Find(&advertisements).Error; err != nil {
		return nil, err
	}
	return advertisements, nil
}

// Example implementation of UpdateAdvertisementPlayStatus
func (am *AdvertisementDBModel) UpdateAdvertisementPlayStatus(advertisementID uint, played bool) error {
	// Replace this example code with your actual logic to update the play status of an advertisement
	if err := am.DB.Model(&Advertisement{}).
		Where("id = ?", advertisementID).
		Update("played", played).
		Error; err != nil {
		return err
	}

	return nil
}

// Example implementation of getNextAdvertisement
func (am *AdvertisementDBModel) GetNextAdvertisement(playlistID uint) (*Advertisement, error) {
	// Replace this example code with your actual logic to get the next advertisement
	var nextAdvertisement Advertisement
	if err := am.DB.Model(&Advertisement{}).
		Where("playlist_id = ? AND start_date <= ? AND end_date >= ? AND played = false", playlistID, time.Now(), time.Now()).
		Order("priority DESC, created_at").
		First(&nextAdvertisement).Error; err != nil {
		return nil, err
	}

	return &nextAdvertisement, nil
}

// GetCurrentAdvertisement retrieves the current advertisement from the database
func (am *AdvertisementDBModel) GetCurrentAdvertisement(id uint) (*Advertisement, error) {
	// Your implementation to retrieve the current advertisement from the database
	var ad Advertisement
	if err := am.DB.First(&ad, id).Error; err != nil {
		return nil, err
	}
	return &ad, nil
}

func (am *AdvertisementDBModel) GetPlaybackRate() (float64, error) {
	// Implement logic to fetch playback rate from the database
	// For example:
	// var ad Advertisement
	// am.DB.First(&ad)
	// return ad.PlaybackRate, nil
	return 1.0, nil // Placeholder value, replace with actual logic
}

// GetVideoIDForAdvertisement fetches the video ID associated with the given advertisement ID
func (m *AdvertisementDBModel) GetVideoIDForAdvertisement(advertisementID uint) (uint, error) {
	var ad Advertisement
	if err := m.DB.First(&ad, advertisementID).Error; err != nil {
		return 0, err
	}
	return ad.VideoID, nil
}

// IncrementPlayCount increments the play count of the advertisement with the given ID
func (am *AdvertisementDBModel) IncrementPlayCount(advertisementID uint) error {
	var advertisement Advertisement
	if err := am.DB.First(&advertisement, advertisementID).Error; err != nil {
		return err
	}

	advertisement.PlayCount++
	return am.DB.Save(&advertisement).Error
}
