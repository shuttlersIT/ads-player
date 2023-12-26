// backend/services/playlist_service.go

package services

import (
	"github.com/shuttlersit/ads-player/backend/models"
)

// PlaylistService provides methods for controlling the playlist
type PlaylistService interface {
	PlayNextVideo() error
	PausePlayback() error
	GetAllPlaylists() ([]models.Playlist, error)
	GetPlaylistByID(id string) (*models.Playlist, error)
	CreatePlaylist(playlist *models.Playlist) (*models.Playlist, error)
	// Add more methods as needed
}

// DefaultPlaylistService is the default implementation of PlaylistService
type DefaultPlaylistService struct {
	playlistModel *models.PlaylistDBModel
	// Add any dependencies or data needed for the service
}

// NewDefaultPlaylistService creates a new DefaultPlaylistService.
func NewDefaultPlaylistService(playlistModel *models.PlaylistDBModel) *DefaultPlaylistService {
	return &DefaultPlaylistService{
		playlistModel: playlistModel,
	}
}

// PlayNextVideo plays the next video in the playlist
func (ps *DefaultPlaylistService) PlayNextVideo() error {
	// Implement logic to play the next video
	return nil
}

// PausePlayback pauses the current playback
func (ps *DefaultPlaylistService) PausePlayback() error {
	// Implement logic to pause playback
	return nil
}

// GetAllPlaylists retrieves all playlists.
func (ps *DefaultPlaylistService) GetAllPlaylists() ([]models.Playlist, error) {
	playlists, err := ps.playlistModel.GetAllPlaylists()
	if err != nil {
		return nil, err
	}
	return playlists, nil
}

// GetPlaylistByID retrieves a playlist by ID.
func (ps *DefaultPlaylistService) GetPlaylistByID(id uint) (*models.Playlist, error) {
	playlist, err := ps.playlistModel.GetPlaylistByID(id)
	if err != nil {
		return nil, err
	}
	return playlist, nil
}

// CreatePlaylist creates a new playlist.
func (ps *DefaultPlaylistService) CreatePlaylist(playlist *models.Playlist) (*models.Playlist, error) {
	err := ps.playlistModel.CreatePlaylist(playlist)
	if err != nil {
		return nil, err
	}
	return playlist, nil
}
