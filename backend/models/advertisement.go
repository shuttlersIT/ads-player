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
	ContentURL            string                 `json:"contentURL"`
	Title                 string                 `json:"title"`
	Description           string                 `json:"description"`
	Duration              int                    `json:"duration"` // Duration in seconds
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
	UpdatedAt             time.Time              `json:"updatedAt"`
	TargetAudience        string                 `json:"targetAudience"`
	MatureContent         bool                   `json:"matureContent"`
	ThumbnailURL          string                 `json:"thumbnailURL"`
	ExternalLinks         []ExternalLink         `json:"externalLinks" gorm:"foreignKey:AdvertisementID"`
	MediaAttachments      []MediaAttachment      `json:"mediaAttachments" gorm:"foreignKey:AdvertisementID"`
	AdvertisementHashtags []AdvertisementHashtag `json:"hashtags" gorm:"foreignKey:AdvertisementID"`
	StartDate             time.Time              `json:"start_date"`
	EndDate               time.Time              `json:"end_date"`
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

// PrivacySetting struct for advertisement privacy settings
type AdvertisementPrivacySetting struct {
	IsPublic       bool   `json:"isPublic" gorm:"default:true"`
	Password       string `json:"password,omitempty"`
	AllowComments  bool   `json:"allowComments" gorm:"default:true"`
	AllowDownloads bool   `json:"allowDownloads" gorm:"default:true"`
	AllowEmbedding bool   `json:"allowEmbedding" gorm:"default:true"`
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

// Hashtag represents a hashtag entity
type Hashtag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Tag       string    `json:"tag"`
	Tags      string    `json:"tags"` // Added Tags field
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RelatedAd struct for storing related advertisements
type RelatedAd struct {
	gorm.Model
	AdvertisementID        uint `json:"-"`
	RelatedAdvertisementID uint `json:"relatedAdvertisementID"`
	Order                  int  `json:"order" gorm:"default:0"`
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
	CreateAdvertisement(ad *Advertisement) error
	UpdateAdvertisement(ad *Advertisement) error
	DeleteAdvertisement(id uint) error
}

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
func (am *AdvertisementDBModel) GetAdvertisementByID(advertisementID uint) (*Advertisement, error) {
	var advertisement Advertisement
	if err := am.DB.Preload("Playlist").Preload("Comments").Preload("Followers").Preload("Contributors").Preload("RelatedAds").First(&advertisement, advertisementID).Error; err != nil {
		return nil, err
	}
	return &advertisement, nil
}

// GetAllAdvertisements fetches all advertisements
func (am *AdvertisementDBModel) GetAllAdvertisements() ([]Advertisement, error) {
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
func (am *AdvertisementDBModel) UpdateAdvertisement(advertisement *Advertisement) error {
	if err := am.DB.Save(advertisement).Error; err != nil {
		return err
	}
	return nil
}

// DeleteAdvertisement deletes an advertisement by its ID
func (am *AdvertisementDBModel) DeleteAdvertisement(advertisementID uint) error {
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
