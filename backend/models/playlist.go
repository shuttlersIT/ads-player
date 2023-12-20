// backend/models/playlist.go

package models

import (
	"errors"
	"log"
	"sort"
	"time"

	"gorm.io/gorm"
)

// Playlist model
type Playlist struct {
	gorm.Model
	ID                           uint              `json:"id" gorm:"primaryKey"`
	Name                         string            `json:"name"`
	CreatorUserID                uint              `json:"creator_user_id"`
	StartTime                    time.Time         `json:"start_time"`
	CreatedAt                    time.Time         `json:"created_at"`
	UpdatedAt                    time.Time         `json:"updated_at"`
	Title                        string            `json:"title"`
	Description                  string            `json:"description"`
	Videos                       []Video           `gorm:"foreignKey:PlaylistID"`
	ChannelID                    uint              `json:"-"`
	Channel                      Channel           `json:"channel"`
	IsFeatured                   bool              `json:"isFeatured" gorm:"default:false"`
	IsPublic                     bool              `json:"isPublic" gorm:"default:true"`
	FeaturedArtwork              string            `json:"featuredArtwork"`
	Tags                         []string          `json:"tags" gorm:"type:varchar(255)[]"`
	IsPlayable                   bool              `json:"isPlayable" gorm:"default:true"`
	PlayCount                    uint              `json:"playCount" gorm:"default:0"`
	LikeCount                    uint              `json:"likeCount" gorm:"default:0"`
	DislikeCount                 uint              `json:"dislikeCount" gorm:"default:0"`
	Comments                     []Comment         `gorm:"foreignKey:PlaylistID"`
	ShareCount                   uint              `json:"shareCount" gorm:"default:0"`
	Advertisements               []Advertisement   `json:"advertisements" gorm:"foreignKey:PlaylistID"`
	Followers                    []User            `gorm:"many2many:user_playlist_followers;"`
	Contributors                 []User            `gorm:"many2many:user_playlist_contributors;"`
	RelatedPlaylists             []RelatedPlaylist `json:"relatedPlaylists" gorm:"foreignKey:PlaylistID"`
	TotalDuration                time.Duration     `json:"totalDuration" gorm:"default:0"`
	LastModified                 int               `json:"lastModified" gorm:"autoUpdateTime"`
	LastAdvertisementScheduledAt time.Time         `json:"lastAdvertisementScheduledAt" gorm:"default:null"`
	PrivacySetting               PrivacySetting    `json:"privacySetting" gorm:"embedded"`
	Location                     Location          `json:"location" gorm:"embedded"`
	CurrentVideo                 *Video            `json:"currentVideo" gorm:"embedded"`
	NextVideo                    *Video            `json:"nextVideo" gorm:"embedded"`
	LastScheduledAt              time.Time         `json:"last_scheduled_at"`
}

// UserPlaylistFollowers model for many-to-many relationship between users and playlists
type UserPlaylistFollowers struct {
	gorm.Model
	UserID     uint
	PlaylistID uint
}

// UserPlaylistContributors model for many-to-many relationship between users and playlists
type UserPlaylistContributors struct {
	gorm.Model
	UserID     uint
	PlaylistID uint
}

// RelatedPlaylist model for representing related playlists
type RelatedPlaylist struct {
	gorm.Model
	PlaylistID        uint
	RelatedPlaylistID uint
}

// PrivacySetting model for defining playlist privacy settings
type PrivacySetting struct {
	IsCollaborative bool `json:"isCollaborative" gorm:"default:false"`
	AllowComments   bool `json:"allowComments" gorm:"default:true"`
	// Add more privacy settings as needed
}

// PlaylistModel handles database operations for Playlist
type PlaylistDBModel struct {
	DB *gorm.DB
}

// NewPlaylistModel creates a new instance of PlaylistModel
func NewPlaylistModel(db *gorm.DB) *PlaylistDBModel {
	return &PlaylistDBModel{
		DB: db,
	}
}

type PlaylistModel interface {
	GetCurrentPlaylist(id uint) (*Playlist, error)
	UpdateLastScheduledTime(id uint, lastScheduledTime time.Time) error
	GetPlaylistsForAdvertisements() ([]Playlist, error)
	GetPlaylistsForAdvertisementsFreshness() ([]Playlist, error)
	GetPlaylistsForAdvertisementsByPopularity() ([]Playlist, error)
	hasActiveAdvertisements() bool
	sortPlaylistsByPopularity([]*Playlist)
	sortPlaylistsByFreshness([]*Playlist)
	calculateTotalViews(playlist *Playlist) int
	hasHighPopularity(playlist *Playlist) bool
	isFreshPlaylist(playlist *Playlist) bool
	GetTotalDuration(playlistID uint) (time.Duration, error)
	GetCurrentVideo(playlistID uint) (*Video, error)
	GetNextVideo(playlistID uint) (*Video, error)
	GetCurrentPlaylistID() (uint, error)
}

// GetCurrentPlaylistID retrieves the ID of the current playlist
func (p *PlaylistDBModel) GetCurrentPlaylistID() uint {
	// Implement the logic to get the current playlist ID
	// For example, you might fetch it from a configuration or database
	return 1 // Replace with your actual implementation
}

// GetCurrentPlaylist retrieves the current playlist from the database
func (am *PlaylistDBModel) GetCurrentPlaylist(id uint) (*Playlist, error) {
	// Your implementation to retrieve the current playlist from the database
	var playlist Playlist
	if err := am.DB.First(&playlist, id).Error; err != nil {
		return nil, err
	}
	return &playlist, nil
}

// GetPlaylistsForAdvertisements fetches playlists that have associated advertisements
func (pm *PlaylistDBModel) GetPlaylistsForAdvertisements() ([]Playlist, error) {
	var playlists []Playlist
	if err := pm.DB.Preload("Advertisements").Find(&playlists).Error; err != nil {
		log.Printf("Error fetching playlists with advertisements: %v", err)
		return nil, errors.New("failed to fetch playlists")
	}

	if len(playlists) == 0 {
		return nil, errors.New("no playlists found with advertisements")
	}

	// Perform additional processing or filtering if needed
	// For example, you can filter out playlists without active advertisements
	// You can also sort playlists based on criteria such as popularity or freshness

	filteredPlaylists := make([]Playlist, 0)
	for _, p := range playlists {
		if hasActiveAdvertisements(&p) && isFreshPlaylist(&p) && hasHighPopularity(&p) {
			filteredPlaylists = append(filteredPlaylists, p)
		}
	}

	if len(filteredPlaylists) == 0 {
		return nil, errors.New("no playlists found with active advertisements, are fresh, and have high popularity")
	}

	// Sort playlists based on a custom ranking algorithm or popularity criteria
	sort.Slice(filteredPlaylists, func(i, j int) bool {
		return calculateTotalViews(&filteredPlaylists[i]) > calculateTotalViews(&filteredPlaylists[j])
	})

	return filteredPlaylists, nil
}

// GetPlaylistsForAdvertisements fetches playlists that have associated advertisements
func (pm *PlaylistDBModel) GetPlaylistsForAdvertisementsFreshness() ([]Playlist, error) {
	var playlists []Playlist
	if err := pm.DB.Preload("Advertisements").Find(&playlists).Error; err != nil {
		log.Printf("Error fetching playlists with advertisements: %v", err)
		return nil, errors.New("failed to fetch playlists")
	}

	if len(playlists) == 0 {
		return nil, errors.New("no playlists found with advertisements")
	}

	// Perform additional processing or filtering if needed
	// For example, you can filter out playlists without active advertisements
	// You can also sort playlists based on criteria such as popularity or freshness

	filteredPlaylists := make([]Playlist, 0)
	for _, p := range playlists {
		if hasActiveAdvertisements(&p) && isFreshPlaylist(&p) {
			filteredPlaylists = append(filteredPlaylists, p)
		}
	}

	if len(filteredPlaylists) == 0 {
		return nil, errors.New("no playlists found with active advertisements and are fresh")
	}

	// Sort playlists based on freshness or popularity
	sortPlaylistsByFreshness(filteredPlaylists)

	return filteredPlaylists, nil
}

/* ][][][][][][][][][][][][][][][][][][][][][][][][][][[][][]][][][][][][][][][][][*/
func (pm *PlaylistDBModel) GetPlaylistsForAdvertisementsByPopularity() ([]Playlist, error) {
	var playlists []Playlist
	if err := pm.DB.Preload("Advertisements").Find(&playlists).Error; err != nil {
		log.Printf("Error fetching playlists with advertisements: %v", err)
		return nil, errors.New("failed to fetch playlists")
	}

	if len(playlists) == 0 {
		return nil, errors.New("no playlists found with advertisements")
	}

	// Perform additional processing or filtering if needed
	// For example, you can filter out playlists without active advertisements
	// You can also sort playlists based on criteria such as popularity or freshness

	filteredPlaylists := make([]Playlist, 0)
	for _, p := range playlists {
		if hasActiveAdvertisements(&p) && isFreshPlaylist(&p) && hasHighPopularity(&p) {
			filteredPlaylists = append(filteredPlaylists, p)
		}
	}

	if len(filteredPlaylists) == 0 {
		return nil, errors.New("no playlists found with active advertisements, are fresh, and have high popularity")
	}

	// Sort playlists based on a custom ranking algorithm or popularity criteria
	sortPlaylistsByPopularity(filteredPlaylists)

	return filteredPlaylists, nil
}

// hasActiveAdvertisements checks if a playlist has active advertisements
func hasActiveAdvertisements(playlist *Playlist) bool {
	for _, ad := range playlist.Advertisements {
		// You can customize the criteria for active advertisements based on your requirements
		if !ad.Played {
			return true
		}
	}
	return false
}

// isFreshPlaylist checks if a playlist is considered fresh (e.g., created within the last 7 days)
func isFreshPlaylist(playlist *Playlist) bool {
	const freshnessThreshold = 7 // Number of days considered as fresh

	return time.Since(playlist.CreatedAt).Hours()/24 <= freshnessThreshold
}

// hasHighPopularity checks if a playlist has high popularity based on custom criteria
func hasHighPopularity(playlist *Playlist) bool {
	// You can implement custom criteria to determine playlist popularity
	// For simplicity, this example considers a playlist with more than 100 views as popular
	return calculateTotalViews(playlist) > 100
}

// calculateTotalViews calculates the total number of views for all advertisements in a playlist
func calculateTotalViews(playlist *Playlist) int {
	totalViews := 0
	for _, ad := range playlist.Advertisements {
		totalViews += int(ad.Analytics.Views)
	}
	return totalViews
}

// sortPlaylistsByFreshness sorts playlists based on freshness or popularity
func sortPlaylistsByFreshness(playlists []Playlist) {
	// You can implement custom sorting logic here
	// For example, sorting by the number of active advertisements or other criteria
	// For simplicity, this example sorts by the playlist's creation time (freshness)
	sort.Slice(playlists, func(i, j int) bool {
		return playlists[i].CreatedAt.After(playlists[j].CreatedAt)
	})
}

// sortPlaylistsByPopularity sorts playlists based on a custom ranking algorithm or popularity criteria
func sortPlaylistsByPopularity(playlists []Playlist) {
	// You can implement a custom sorting algorithm based on popularity or other criteria
	// For simplicity, this example sorts playlists based on the total number of views
	sort.Slice(playlists, func(i, j int) bool {
		return calculateTotalViews(&playlists[i]) > calculateTotalViews(&playlists[j])
	})
}

func (pm *PlaylistDBModel) GetCurrentVideo(playlistID uint) (*Video, error) {
	// Implement logic to fetch information about the currently playing video from the database for the specified playlistID
	// For example:
	// var v Video
	// pm.DB.Where("playlist_id = ? AND playing = ?", playlistID, true).First(&v)
	// return &v, nil
	return &Video{ID: 1, ThumbnailURL: "current_thumbnail.jpg"}, nil // Placeholder value, replace with actual logic
}

func (pm *PlaylistDBModel) GetNextVideo(playlistID uint) (*Video, error) {
	// Implement logic to fetch information about the next video in the playlist from the database for the specified playlistID
	// For example:
	// var v Video
	// pm.DB.Where("playlist_id = ? AND playing = ?", playlistID, false).First(&v)
	// return &v, nil
	return &Video{ID: 2, ThumbnailURL: "next_thumbnail.jpg"}, nil // Placeholder value, replace with actual logic
}

// GetTotalDuration retrieves the total duration of the current playlist
func (m *PlaylistDBModel) GetTotalDuration(playlistID uint) (time.Duration, error) {
	var playlist Playlist
	if err := m.DB.First(&playlist, playlistID).Error; err != nil {
		return 0, err
	}
	return playlist.TotalDuration, nil
}

// UpdateLastScheduledTime updates the last scheduled time for a playlist
func (m *PlaylistDBModel) UpdateLastScheduledTime(playlistID uint, lastScheduledAt time.Time) error {
	var playlist Playlist
	if err := m.DB.First(&playlist, playlistID).Error; err != nil {
		return err
	}
	playlist.LastScheduledAt = lastScheduledAt
	return m.DB.Save(&playlist).Error
}

// UpdateLastScheduledTime updates the last scheduled time for a playlist
func (pm *PlaylistDBModel) UpdateLastAdvertisementScheduledTime(playlistID uint, lastScheduledTime time.Time) error {
	var playlist Playlist
	if err := pm.DB.First(&playlist, playlistID).Error; err != nil {
		return err
	}

	playlist.LastAdvertisementScheduledAt = lastScheduledTime
	if err := pm.DB.Save(&playlist).Error; err != nil {
		return err
	}

	return nil
}

// GetStartTime gets the start time of the playlist with the given ID
func (m *PlaylistDBModel) GetStartTime(playlistID uint) (time.Time, error) {
	// Replace with the actual implementation to get the start time of the playlist
	// For example, query the database, fetch from cache, etc.
	// This is just a simulated example

	// Assume there is a field in the Playlist model representing the start time
	var playlist Playlist
	if err := m.DB.Where("id = ?", playlistID).First(&playlist).Error; err != nil {
		return time.Time{}, err
	}

	return playlist.StartTime, nil
}
