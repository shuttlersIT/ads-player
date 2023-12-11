// backend/models/advertisement_play_event.go

package models

import (
	"time"
)

// AdvertisementPlayEvent represents an event when an advertisement is played
type AdvertisementPlayEvent struct {
	ID              uint `gorm:"primaryKey"`
	AdvertisementID uint
	PlaylistID      uint
	PlayTime        time.Time
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
