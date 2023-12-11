// backend/models/advertisement.go

package models

import (
	"time"

	"gorm.io/gorm"
)

// Advertisement model
type Advertisement struct {
	gorm.Model
	PlaylistID       uint                   `json:"-"`
	Playlist         Playlist               `json:"playlist"`
	ContentURL       string                 `json:"contentURL"`
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	Duration         int                    `json:"duration"` // Duration in seconds
	ScheduledAt      time.Time              `json:"scheduledAt"`
	Played           bool                   `json:"played" gorm:"default:false"`
	ClickThroughURL  string                 `json:"clickThroughURL"`
	Analytics        AdvertisementAnalytics `json:"analytics" gorm:"embedded"`
	IsFeatured       bool                   `json:"isFeatured" gorm:"default:false"`
	IsPublic         bool                   `json:"isPublic" gorm:"default:true"`
	Tags             []string               `json:"tags" gorm:"type:varchar(255)[]"`
	LikeCount        uint                   `json:"likeCount" gorm:"default:0"`
	DislikeCount     uint                   `json:"dislikeCount" gorm:"default:0"`
	Comments         []Comment              `gorm:"foreignKey:AdvertisementID"`
	ShareCount       uint                   `json:"shareCount" gorm:"default:0"`
	Followers        []User                 `gorm:"many2many:user_advertisement_followers;"`
	Contributors     []User                 `gorm:"many2many:user_advertisement_contributors;"`
	RelatedAds       []RelatedAd            `json:"relatedAds" gorm:"foreignKey:AdvertisementID"`
	LastModified     int                    `json:"lastModified" gorm:"autoUpdateTime"`
	PrivacySetting   PrivacySetting         `json:"privacySetting" gorm:"embedded"`
	Location         Location               `json:"location" gorm:"embedded"`
	VideoQuality     string                 `json:"videoQuality"`
	AudioQuality     string                 `json:"audioQuality"`
	Caption          string                 `json:"caption"`
	Language         string                 `json:"language"`
	TargetAudience   string                 `json:"targetAudience"`
	MatureContent    bool                   `json:"matureContent"`
	ThumbnailURL     string                 `json:"thumbnailURL"`
	ExternalLinks    []ExternalLink         `json:"externalLinks" gorm:"foreignKey:AdvertisementID"`
	MediaAttachments []MediaAttachment      `json:"mediaAttachments" gorm:"foreignKey:AdvertisementID"`
	Hashtags         []Hashtag              `json:"hashtags" gorm:"foreignKey:AdvertisementID"`
}

// AdvertisementAnalytics struct for tracking advertisement analytics
type AdvertisementAnalytics struct {
	Views                uint    `json:"views" gorm:"default:0"`
	Clicks               uint    `json:"clicks" gorm:"default:0"`
	Comments             uint    `json:"comments" gorm:"default:0"`
	Shares               uint    `json:"shares" gorm:"default:0"`
	Likes                uint    `json:"likes" gorm:"default:0"`
	Dislikes             uint    `json:"dislikes" gorm:"default:0"`
	Followers            uint    `json:"followers" gorm:"default:0"`
	PlayCount            uint    `json:"playCount" gorm:"default:0"`
	TotalDurationWatched int     `json:"totalDurationWatched" gorm:"default:0"`
	ClickThroughCount    uint    `json:"clickThroughCount" gorm:"default:0"`
	ConversionRate       float64 `json:"conversionRate" gorm:"default:0.0"`
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
	AdvertisementID uint   `json:"-"`
	Text            string `json:"text"`
	// Add more fields as needed
}

// RelatedAd struct for storing related advertisements
type RelatedAd struct {
	gorm.Model
	AdvertisementID        uint `json:"-"`
	RelatedAdvertisementID uint `json:"relatedAdvertisementID"`
	Order                  int  `json:"order" gorm:"default:0"`
}

// AdvertisementModel handles database operations for Advertisement
type AdvertisementModel struct {
	DB *gorm.DB
}

// NewAdvertisementModel creates a new instance of AdvertisementModel
func NewAdvertisementModel(db *gorm.DB) *AdvertisementModel {
	return &AdvertisementModel{
		DB: db,
	}
}

// GetNextAdvertisementForPlaylist fetches the next advertisement to play for a playlist
func (am *AdvertisementModel) GetNextAdvertisementForPlaylist(playlistID uint) (*Advertisement, error) {
	var advertisement Advertisement
	if err := am.DB.Where("playlist_id = ? AND scheduled_at <= ? AND played = ?", playlistID, time.Now(), false).Order("scheduled_at").First(&advertisement).Error; err != nil {
		return nil, err
	}

	return &advertisement, nil
}

// MarkAdvertisementAsPlayed marks an advertisement as played (update status or other relevant fields)
func (am *AdvertisementModel) MarkAdvertisementAsPlayed(advertisementID uint) error {
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
func (am *AdvertisementModel) GetAdvertisementByID(advertisementID uint) (*Advertisement, error) {
	var advertisement Advertisement
	if err := am.DB.Preload("Playlist").Preload("Comments").Preload("Followers").Preload("Contributors").Preload("RelatedAds").First(&advertisement, advertisementID).Error; err != nil {
		return nil, err
	}
	return &advertisement, nil
}

// GetAllAdvertisements fetches all advertisements
func (am *AdvertisementModel) GetAllAdvertisements() ([]Advertisement, error) {
	var advertisements []Advertisement
	if err := am.DB.Preload("Playlist").Preload("Comments").Preload("Followers").Preload("Contributors").Preload("RelatedAds").Find(&advertisements).Error; err != nil {
		return nil, err
	}
	return advertisements, nil
}

// CreateAdvertisement creates a new advertisement
func (am *AdvertisementModel) CreateAdvertisement(advertisement *Advertisement) error {
	if err := am.DB.Create(advertisement).Error; err != nil {
		return err
	}
	return nil
}

// UpdateAdvertisement updates an existing advertisement
func (am *AdvertisementModel) UpdateAdvertisement(advertisement *Advertisement) error {
	if err := am.DB.Save(advertisement).Error; err != nil {
		return err
	}
	return nil
}

// DeleteAdvertisement deletes an advertisement by its ID
func (am *AdvertisementModel) DeleteAdvertisement(advertisementID uint) error {
	if err := am.DB.Delete(&Advertisement{}, advertisementID).Error; err != nil {
		return err
	}
	return nil
}

// GetAdvertisementsByPlaylistID fetches all advertisements for a specific playlist
func (am *AdvertisementModel) GetAdvertisementsByPlaylistID(playlistID uint) ([]Advertisement, error) {
	var advertisements []Advertisement
	if err := am.DB.Where("playlist_id = ?", playlistID).Preload("Playlist").Preload("Comments").Preload("Followers").Preload("Contributors").Preload("RelatedAds").Find(&advertisements).Error; err != nil {
		return nil, err
	}
	return advertisements, nil
}
