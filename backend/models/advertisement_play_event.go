// backend/models/advertisement_play_event.go

package models

import (
	"time"

	"gorm.io/gorm"
)

// AdvertisementPlayEvent represents an event when an advertisement is played
type AdvertisementPlayEvent struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	AdvertisementID uint      `json:"advertisement_id"`
	PlaylistID      uint      `json:"playlist_id"`
	ChannelID       uint      `gorm:"not null" json:"channel_id"`
	VideoID         uint      `gorm:"not null" json:"video_id"`
	UserID          uint      `json:"user_id"`
	PlayTime        time.Time `json:"play_time"`
	PlayedAt        time.Time `gorm:"not null" json:"playedAt"`
	DurationWatched uint      `gorm:"not null" json:"durationWatched"`
	CreatedAt       time.Time `json:"created_at"`
}

// TableName sets the table name for the AdvertisementPlayEvent model.
func (AdvertisementPlayEvent) TableName() string {
	return "advertisement_play_events"
}

// LogAdvertisementPlayEvent logs an event when an advertisement is played
func (am *AdvertisementDBModel) LogAdvertisementPlayEvent(advertisementID, playlistID uint) error {
	playEvent := AdvertisementPlayEvent{
		AdvertisementID: advertisementID,
		PlaylistID:      playlistID,
		PlayTime:        time.Now(),
	}

	if err := am.DB.Create(&playEvent).Error; err != nil {
		return err
	}

	return nil
}

// AdvertisementPlayEventDBModel provides methods for database operations related to AdvertisementPlayEvent.
type AdvertisementPlayEventDBModel struct {
	DB *gorm.DB
}

// NewAdvertisementPlayEventDBModel creates a new instance of AdvertisementPlayEventDBModel.
func NewAdvertisementPlayEventDBModel(db *gorm.DB) *AdvertisementPlayEventDBModel {
	return &AdvertisementPlayEventDBModel{DB: db}
}

// CreateAdvertisementPlayEvent creates a new advertisement play event in the database.
func (m *AdvertisementPlayEventDBModel) CreateAdvertisementPlayEvent(event *AdvertisementPlayEvent) error {
	return m.DB.Create(event).Error
}

// SaveAdvertisementPlayEvent saves an advertisement play event in the database.
func (m *AdvertisementPlayEventDBModel) SaveAdvertisementPlayEvent(event *AdvertisementPlayEvent) error {
	return m.DB.Save(event).Error
}

// GetAdvertisementPlayEvents retrieves the play events for a specific advertisement.
func (m *AdvertisementPlayEventDBModel) GetAdvertisementPlayEvents(advertisementID uint) ([]AdvertisementPlayEvent, error) {
	var events []AdvertisementPlayEvent
	err := m.DB.Where("advertisement_id = ?", advertisementID).Find(&events).Error
	return events, err
}

// GetPlaylistAdvertisementPlayEvents retrieves play events for advertisements in a playlist.
func (m *AdvertisementPlayEventDBModel) GetPlaylistAdvertisementPlayEvents(playlistID uint) ([]AdvertisementPlayEvent, error) {
	var events []AdvertisementPlayEvent
	err := m.DB.Where("playlist_id = ?", playlistID).Find(&events).Error
	return events, err
}
