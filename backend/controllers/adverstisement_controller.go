// backend/controllers/advertisement_controller.go

package controllers

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/shuttlersit/ads-player/backend/models"
)

// AdvertisementController handles advertisement scheduling and playback
type AdvertisementController struct {
	PlaylistModel      *models.PlaylistModel
	AdvertisementModel *models.AdvertisementModel
	PlaybackService    PlaybackService
}

// NewAdvertisementController creates a new instance of AdvertisementController
func NewAdvertisementController(playlistModel *models.PlaylistModel, advertisementModel *models.AdvertisementModel, playbackService PlaybackService) *AdvertisementController {
	return &AdvertisementController{
		PlaylistModel:      playlistModel,
		AdvertisementModel: advertisementModel,
		PlaybackService:    playbackService,
	}
}

// ScheduleAdvertisements sets up a cron job to schedule advertisements
func (ac *AdvertisementController) ScheduleAdvertisements() error {
	c := cron.New()

	_, err := c.AddFunc("*/5 * * * *", func() {
		// Get playlists that are eligible for advertisements
		playlists, err := ac.PlaylistModel.GetPlaylistsForAdvertisements()
		if err != nil {
			fmt.Println("Error fetching playlists for advertisements:", err)
			return
		}

		// Loop through playlists and schedule advertisements
		for _, playlist := range playlists {
			err := ac.ScheduleAdvertisementForPlaylist(playlist)
			if err != nil {
				fmt.Printf("Error scheduling advertisement for playlist %d: %v\n", playlist.ID, err)
			}
		}
	})
	if err != nil {
		return err
	}

	// Start the cron scheduler
	c.Start()

	return nil
}

// ScheduleAdvertisementForPlaylist schedules an advertisement for a specific playlist
func (ac *AdvertisementController) ScheduleAdvertisementForPlaylist(playlist models.Playlist) error {
	// Get the next advertisement to play for the playlist
	advertisement, err := ac.AdvertisementModel.GetNextAdvertisementForPlaylist(playlist.ID)
	if err != nil {
		return err
	}

	// Check if there is an advertisement to play
	if advertisement != nil {
		// Perform the logic to play the advertisement
		err := ac.PlayAdvertisement(advertisement, playlist)
		if err != nil {
			return err
		}

		// Update the last scheduled time for the playlist
		err = ac.PlaylistModel.UpdateLastScheduledTime(playlist.ID, time.Now())
		if err != nil {
			return err
		}
	}

	return nil
}

// PlayAdvertisement plays an advertisement for a playlist
func (ac *AdvertisementController) PlayAdvertisement(advertisement *models.Advertisement, playlist models.Playlist) error {
	fmt.Printf("Playing advertisement %d for playlist %d\n", advertisement.ID, playlist.ID)

	// Use the playback service to play the advertisement
	err := ac.PlaybackService.Play(advertisement)
	if err != nil {
		fmt.Printf("Error playing advertisement %d: %v\n", advertisement.ID, err)
		return err
	}

	// Update the advertisement's play status
	err = ac.UpdateAdvertisementPlayStatus(advertisement.ID, true)
	if err != nil {
		fmt.Printf("Error updating play status for advertisement %d: %v\n", advertisement.ID, err)
		return err
	}

	// Log the play event
	err = ac.LogAdvertisementPlayEvent(advertisement.ID, playlist.ID)
	if err != nil {
		fmt.Printf("Error logging play event for advertisement %d: %v\n", advertisement.ID, err)
		return err
	}

	return nil
}

// UpdateAdvertisementPlayStatus updates the play status of an advertisement
func (ac *AdvertisementController) UpdateAdvertisementPlayStatus(advertisementID uint, played bool) error {
	// Replace the placeholder code with your actual database update operation
	advertisement, err := ac.AdvertisementModel.GetAdvertisementByID(advertisementID)
	if err != nil {
		return err
	}

	advertisement.Played = played

	err = ac.AdvertisementModel.UpdateAdvertisement(advertisement)
	if err != nil {
		return err
	}

	return nil
}

// LogAdvertisementPlayEvent logs the play event of an advertisement for a playlist
func (ac *AdvertisementController) LogAdvertisementPlayEvent(advertisementID, playlistID uint) error {
	// Replace the placeholder code with your actual database log operation
	playEvent := models.AdvertisementPlayEvent{
		AdvertisementID: advertisementID,
		PlaylistID:      playlistID,
		PlayTime:        time.Now(),
	}

	err := ac.AdvertisementModel.LogAdvertisementPlayEvent(playEvent.AdvertisementID, playEvent.PlaylistID)
	if err != nil {
		return err
	}

	return nil
}
