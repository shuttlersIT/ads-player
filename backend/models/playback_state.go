// backend/models/playback_state.go

package models

import (
	"time"
)

// PlaybackState represents the current playback state
type PlaybackState struct {
	Status               string         `json:"status"`
	CurrentAdvertisement *Advertisement `json:"current_ad"`
	CurrentPlaylist      *Playlist      `json:"current_playlist"`
	ElapsedTime          time.Duration  `json:"elapsed_time"`
	NextAdvertisement    *Advertisement `json:"next_ad"`

	// Additional fields for a more detailed playback state
	PlaybackRate     float64       `json:"playback_rate"`     // Playback speed (e.g., 1.0 for normal speed)
	TotalDuration    time.Duration `json:"total_duration"`    // Total duration of the current playlist or content
	RemainingTime    time.Duration `json:"remaining_time"`    // Remaining time of the current playlist or content
	CurrentVideo     *Video        `json:"current_video"`     // Details about the current video being played
	NextVideo        *Video        `json:"next_video"`        // Details about the next video to be played
	CurrentTimestamp time.Time     `json:"current_timestamp"` // Timestamp of the current playback position
	ThumbnailURL     string        `json:"thumbnail_url"`     // URL of the video thumbnail
	PlaybackHistory  []*Video      `json:"playback_history"`  // History of videos played in the current session

	// Add more relevant information as needed
}

// NewPlaybackState creates a new instance of PlaybackState
func NewPlaybackState(status string, currentAd *Advertisement, currentPlaylist *Playlist, elapsedTime time.Duration, nextAd *Advertisement, playbackRate float64, totalDuration time.Duration, remainingTime time.Duration, currentVideo *Video, nextVideo *Video, currentTimestamp time.Time, thumbnailURL string, playbackHistory []*Video) *PlaybackState {
	return &PlaybackState{
		Status:               status,
		CurrentAdvertisement: currentAd,
		CurrentPlaylist:      currentPlaylist,
		ElapsedTime:          elapsedTime,
		NextAdvertisement:    nextAd,
		PlaybackRate:         playbackRate,
		TotalDuration:        totalDuration,
		RemainingTime:        remainingTime,
		CurrentVideo:         currentVideo,
		NextVideo:            nextVideo,
		CurrentTimestamp:     currentTimestamp,
		ThumbnailURL:         thumbnailURL,
		PlaybackHistory:      playbackHistory,
		// Initialize additional fields as needed
	}
}
