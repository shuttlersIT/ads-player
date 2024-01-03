// backend/models/playback_state.go

package models

import (
	"time"

	"gorm.io/gorm"
)

// PlaybackState represents the current playback state
type PlaybackState struct {
	ID                   uint           `gorm:"primaryKey" json:"id"`
	VideoID              uint           `gorm:"not null" json:"video_id"`
	Position             int            `json:"position"`
	IsPlaying            bool           `gorm:"not null" json:"is_playing"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            time.Time      `gorm:"index" json:"-"`
	Status               string         `json:"status"`
	CurrentAdvertisement *Advertisement `json:"current_ad"`
	CurrentPlaylist      *Playlist      `json:"current_playlist"`
	ElapsedTime          time.Duration  `json:"elapsed_time"`
	NextAdvertisement    *Advertisement `json:"next_ad"`
	PlaylistID           uint           `gorm:"not null" json:"playlistId"`
	ChannelID            uint           `gorm:"not null" json:"channelId"`
	PlaybackRate         float64        `json:"playback_rate"`                     // Playback speed (e.g., 1.0 for normal speed)
	TotalDuration        time.Duration  `json:"total_duration"`                    // Total duration of the current playlist or content
	RemainingTime        time.Duration  `json:"remaining_time"`                    // Remaining time of the current playlist or content
	CurrentVideo         *Video         `json:"current_video"`                     // Details about the current video being played
	NextVideo            *Video         `json:"next_video"`                        // Details about the next video to be played
	CurrentTimestamp     time.Time      `gorm:"not null" json:"current_timestamp"` // Timestamp of the current playback position
	ThumbnailURL         string         `json:"thumbnail_url"`                     // URL of the video thumbnail
	PlaybackHistory      []*Video       `json:"playback_history"`                  // History of videos played in the current session
	// Additional fields for a more detailed playback state
}

// TableName sets the table name for the PlaybackState model.
func (PlaybackState) TableName() string {
	return "playback_states"
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

// PlaybackStateDBModel provides methods for database operations related to PlaybackState.
type PlaybackStateDBModel struct {
	DB *gorm.DB
}

// NewPlaybackStateDBModel creates a new instance of PlaybackStateDBModel.
func NewPlaybackStateDBModel(db *gorm.DB) *PlaybackStateDBModel {
	return &PlaybackStateDBModel{DB: db}
}

// SavePlaybackState saves the playback state in the database.
func (m *PlaybackStateDBModel) SavePlaybackState(state *PlaybackState) error {
	return m.DB.Save(state).Error
}

// CreatePlaybackState creates a new playback state in the database.
func (m *PlaybackStateDBModel) CreatePlaybackState(state *PlaybackState) error {
	return m.DB.Create(state).Error
}

// GetPlaybackStateByVideoAndPlaylist retrieves the playback state for a video in a playlist from the database.
func (m *PlaybackStateDBModel) GetPlaybackStateByVideoAndPlaylist(videoID, playlistID uint) (*PlaybackState, error) {
	var state PlaybackState
	err := m.DB.Where("video_id = ? AND playlist_id = ?", videoID, playlistID).First(&state).Error
	return &state, err
}

// UpdatePlaybackState updates the playback state in the database.
func (m *PlaybackStateDBModel) UpdatePlaybackState(state *PlaybackState) error {
	return m.DB.Save(state).Error
}

// GetPlaybackState retrieves the playback state for a video by its ID.
func (m *PlaybackStateDBModel) GetPlaybackState(videoID uint) (*PlaybackState, error) {
	var state PlaybackState
	err := m.DB.Where("video_id = ?", videoID).First(&state).Error
	return &state, err
}

// DeletePlaybackState deletes the playback state for a video by its ID.
func (m *PlaybackStateDBModel) DeletePlaybackState(videoID uint) error {
	return m.DB.Delete(&PlaybackState{}, "video_id = ?", videoID).Error
}
