// backend/models/playlist.go

package models

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
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
	CreatedAt                    time.Time         `gorm:"autoCreateTime" json:"creation_date"`
	UpdatedAt                    time.Time         `json:"updated_at"`
	DeletedAt                    time.Time         `gorm:"index" json:"-"`
	Title                        string            `json:"title"`
	Description                  string            `json:"description"`
	Videos                       []Video           `gorm:"many2many:playlist_videos;" json:"videos"`
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
	Advertisements               []Advertisement   `gorm:"many2many:playlist_advertisements;" json:"advertisements"`
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
	IsPlaying                    bool              `json:"isPlaying" gorm:"default:false"`
	Order                        int               `json:"order" gorm:"default:0"`
}

// TableName sets the table name for the Playlist model.
func (Playlist) TableName() string {
	return "playlists"
}

// UserPlaylistFollowers model for many-to-many relationship between users and playlists
type UserPlaylistFollowers struct {
	gorm.Model
	UserID     uint
	PlaylistID uint
}

// TableName sets the table name for the Playlist model.
func (UserPlaylistFollowers) TableName() string {
	return "UserPlaylistsFollowers"
}

// UserPlaylistContributors model for many-to-many relationship between users and playlists
type UserPlaylistContributors struct {
	gorm.Model
	UserID     uint
	PlaylistID uint
}

// TableName sets the table name for the Playlist model.
func (UserPlaylistContributors) TableName() string {
	return "UserPlaylistContributors"
}

// RelatedPlaylist model for representing related playlists
type RelatedPlaylist struct {
	gorm.Model
	PlaylistID        uint
	RelatedPlaylistID uint
}

// TableName sets the table name for the Playlist model.
func (RelatedPlaylist) TableName() string {
	return "RelatedPlaylists"
}

// PrivacySetting model for defining playlist privacy settings
type PrivacySetting struct {
	IsCollaborative bool `json:"isCollaborative" gorm:"default:false"`
	AllowComments   bool `json:"allowComments" gorm:"default:true"`
	// Add more privacy settings as needed
}

// TableName sets the table name for the Playlist model.
func (PrivacySetting) TableName() string {
	return "PrivacySetting"
}

// PlaylistModel handles database operations for Playlist
type PlaylistDBModel struct {
	DB           *gorm.DB
	VideoDBModel *VideoDBModel
}

// NewPlaylistModel creates a new instance of PlaylistModel
func NewPlaylistModel(db *gorm.DB, videoDBModel *VideoDBModel) *PlaylistDBModel {
	return &PlaylistDBModel{
		DB:           db,
		VideoDBModel: videoDBModel,
	}
}

type PlaylistModel interface {
	CreatePlaylist(playlist *Playlist) (*Playlist, error)
	CreatePlaylist2(playlist *Playlist) error
	GetAllPlaylists() ([]Playlist, error)
	GetAllPlaylists2() ([]Playlist, error)
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
	UpdateLastAdvertisementScheduledTime(playlistID uint, lastScheduledTime time.Time) error
	GetStartTime(playlistID uint) (time.Time, error)
	GetContributors(playlistID uint) ([]User, error)
	RemoveContributor(playlistID, userID uint) error
	AddContributor(playlistID, userID uint) error
	GetPlaylistsByUserID(userID uint) ([]Playlist, error)
	RemoveFollower(playlistID, userID uint) error
	AddFollower(playlistID, userID uint) error
	GetFollowers(playlistID uint) ([]User, error)
	RemoveAdvertisementFromPlaylist(playlistID, adID uint) error
	AddAdvertisementToPlaylist(playlistID uint, ad *Advertisement) error
	UpdatePlaylist(playlist *Playlist) (*Playlist, error)
	UpdatePlaylist2(playlist *Playlist) error
	GetPlaylistByID(playlistID uint) (*Playlist, error)
	GetPlaylistByID2(playlistID uint) (*Playlist, error)
	DeletePlaylist(playlistID uint) error
	DeletePlaylist2(playlistID uint) error
	UpdatePlaylistDetails(playlistID uint, updatedDetails Playlist) error
	AddRelatedPlaylist(playlistID, relatedPlaylistID uint) error
	RemoveRelatedPlaylist(playlistID, relatedPlaylistID uint) error
	GetPlaybackStatus(playlistID uint) string
	GetPlaylistInfo(playlistID uint) string
	AddToPlaylist(playlistID uint, videoID uint)
	RemoveFromPlaylist(playlistID uint, videoID uint)
	GetCurrentVideoInfo(playlistID uint) string
	PlayVideo(playlistID uint, videoID uint)
	PauseVideo(playlistID uint)
	ResumeVideo(playlistID uint)
	IsVideoPlaying2() bool
	SkipToPosition(playlistID uint, position int) (int, error)
	IsVideoPlaying(playlistID uint) bool
	AdjustVolume(playlistID uint, delta int) int
	AddVideoToPlaylist(playlistID uint, video *Video) error
	RemoveVideoFromPlaylist(playlistID uint, videoID uint) error
	ShufflePlaylist(playlistID uint) error
	GetCurrentlyPlayingVideoID(playlistID uint) uint
	GetVideoQueue(playlistID uint) ([]Video, error)
}

// CreatePlaylist creates a new playlist in the database.
func (m *PlaylistDBModel) CreatePlaylist(playlist *Playlist) error {
	return m.DB.Create(playlist).Error
}

// GetPlaylistByID retrieves a playlist by its ID.
func (m *PlaylistDBModel) GetPlaylistByID(id uint) (*Playlist, error) {
	var playlist Playlist
	err := m.DB.Preload("Videos").Where("id = ?", id).First(&playlist).Error
	return &playlist, err
}

// UpdatePlaylist updates the details of an existing playlist.
func (m *PlaylistDBModel) UpdatePlaylist(playlist *Playlist) error {
	return m.DB.Save(playlist).Error
}

// DeletePlaylist deletes a playlist from the database.
func (m *PlaylistDBModel) DeletePlaylist(id uint) error {
	return m.DB.Delete(&Playlist{}, id).Error
}

// GetAllPlaylists retrieves all playlists from the database.
func (m *PlaylistDBModel) GetAllPlaylists() ([]Playlist, error) {
	var playlists []Playlist
	err := m.DB.Find(&playlists).Error
	return playlists, err
}

/*-----------------------------------------------------------------------------------------------------------------------*/

func (m *PlaylistDBModel) GetPlaylistByID2(playlistID uint) (*Playlist, error) {
	var playlist Playlist
	if err := m.DB.Preload("Videos").First(&playlist, playlistID).Error; err != nil {
		return nil, err
	}
	return &playlist, nil
}

func (m *PlaylistDBModel) GetAllPlaylists2() ([]*Playlist, error) {
	var playlists []*Playlist
	if err := m.DB.Preload("Videos").Find(&playlists).Error; err != nil {
		return nil, err
	}
	return playlists, nil
}

// GetAllPlaylists retrieves all playlists from the database
func (m *PlaylistDBModel) GetAllPlaylists3() ([]Playlist, error) {
	var playlists []Playlist
	if err := m.DB.Preload("Contributors").Find(&playlists).Error; err != nil {
		return nil, err
	}

	return playlists, nil
}

// CreatePlaylist creates a new playlist
func (m *PlaylistDBModel) CreatePlaylist2(playlist *Playlist) (*Playlist, error) {
	if err := m.DB.Create(playlist).Error; err != nil {
		return nil, err
	}
	return playlist, nil
}

// CreatePlaylist creates a new playlist
func (m *PlaylistDBModel) CreatePlaylist3(playlist *Playlist) error {
	if err := m.DB.Create(playlist).Error; err != nil {
		return err
	}

	return nil
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
	var v Video
	pm.DB.Where("playlist_id = ? AND playing = ?", playlistID, true).First(&v)
	return &v, nil
	//return &Video{ID: 1, ThumbnailURL: "current_thumbnail.jpg"}, nil // Placeholder value, replace with actual logic
}

// GetNextVideo retrieves the next video in the playlist.
func (m *PlaylistDBModel) GetNextVideo(playlistID uint) (*Video, error) {
	// Implement logic to get the next video in the playlist
	// Placeholder code to demonstrate the idea:
	// This could involve checking the order or index to determine the next video.
	playlist, err := m.GetPlaylistByID(playlistID)
	if err != nil || len(playlist.Videos) == 0 {
		return nil, errors.New("no videos in the playlist")
	}

	// Increment the order to get the next video
	playlist.Order++

	// Reset the order to the beginning if it exceeds the length of the playlist
	if playlist.Order >= len(playlist.Videos) {
		playlist.Order = 0
	}

	// Save the changes to the database
	if err := m.DB.Save(playlist).Error; err != nil {
		return nil, err
	}

	return &playlist.Videos[playlist.Order], nil
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

func (m *PlaylistDBModel) UpdatePlaylist2(playlist *Playlist) (*Playlist, error) {
	if err := m.DB.Save(playlist).Error; err != nil {
		return nil, err
	}
	return playlist, nil
}

// UpdatePlaylist updates the information of a playlist in the database
func (m *PlaylistDBModel) UpdatePlaylist3(playlist *Playlist) error {
	return m.DB.Save(playlist).Error
}

// AddAdvertisementToPlaylist adds an advertisement to a playlist
func (m *PlaylistDBModel) AddAdvertisementToPlaylist(playlistID uint, ad *Advertisement) error {
	var playlist Playlist
	if err := m.DB.First(&playlist, playlistID).Error; err != nil {
		return err
	}

	// Assuming you have a relationship set up between Playlist and Advertisement models
	// Adjust the code accordingly based on your database schema
	if err := m.DB.Model(&playlist).Association("Advertisements").Append(ad); err != nil {
		return err
	}

	return nil
}

// RemoveAdvertisementFromPlaylist removes an advertisement from a playlist
func (m *PlaylistDBModel) RemoveAdvertisementFromPlaylist(playlistID, adID uint) error {
	var playlist Playlist
	if err := m.DB.Preload("Advertisements").First(&playlist, playlistID).Error; err != nil {
		return err
	}

	// Assuming you have a relationship set up between Playlist and Advertisement models
	// Adjust the code accordingly based on your database schema
	for i, ad := range playlist.Advertisements {
		if ad.ID == adID {
			playlist.Advertisements = append(playlist.Advertisements[:i], playlist.Advertisements[i+1:]...)
			break
		}
	}

	if err := m.DB.Save(&playlist).Error; err != nil {
		return err
	}

	return nil
}

// GetFollowers retrieves followers of a playlist
func (m *PlaylistDBModel) GetFollowers(playlistID uint) ([]User, error) {
	var playlist Playlist
	if err := m.DB.Preload("Followers").First(&playlist, playlistID).Error; err != nil {
		return nil, err
	}

	// Assuming you have a relationship set up between Playlist and User models for followers
	// Adjust the code accordingly based on your database schema
	var followers []User
	for _, follower := range playlist.Followers {
		followers = append(followers, User{UserID: follower.ID, Username: "Follower"}) // Replace with actual user details
	}

	return followers, nil
}

// AddFollower adds a user as a follower to a playlist
func (m *PlaylistDBModel) AddFollower(playlistID, userID uint) error {
	var playlist Playlist
	if err := m.DB.Preload("Followers").First(&playlist, playlistID).Error; err != nil {
		return err
	}

	// Assuming you have a relationship set up between Playlist and User models for followers
	// Adjust the code accordingly based on your database schema
	if err := m.DB.Model(&playlist).Association("Followers").Append(&UserPlaylistFollowers{UserID: userID}); err != nil {
		return err
	}

	return nil
}

// RemoveFollower removes a follower from a playlist
func (m *PlaylistDBModel) RemoveFollower(playlistID, userID uint) error {
	var playlist Playlist
	if err := m.DB.Preload("Followers").First(&playlist, playlistID).Error; err != nil {
		return err
	}

	// Assuming you have a relationship set up between Playlist and User models for followers
	// Adjust the code accordingly based on your database schema
	for i, follower := range playlist.Followers {
		if follower.UserID == userID {
			playlist.Followers = append(playlist.Followers[:i], playlist.Followers[i+1:]...)
			break
		}
	}

	if err := m.DB.Save(&playlist).Error; err != nil {
		return err
	}

	return nil
}

// GetPlaylistsByUserID retrieves playlists created by a specific user
func (m *PlaylistDBModel) GetPlaylistsByUserID(userID uint) ([]Playlist, error) {
	var playlists []Playlist
	if err := m.DB.Preload("Videos").Preload("Advertisements").Where("creator_user_id = ?", userID).Find(&playlists).Error; err != nil {
		return nil, err
	}
	return playlists, nil
}

// AddContributor adds a user as a contributor to a playlist
func (m *PlaylistDBModel) AddContributor(playlistID, userID uint) error {
	var playlist Playlist
	if err := m.DB.Preload("Contributors").First(&playlist, playlistID).Error; err != nil {
		return err
	}

	// Assuming you have a relationship set up between Playlist and User models for contributors
	// Adjust the code accordingly based on your database schema
	if err := m.DB.Model(&playlist).Association("Contributors").Append(&UserPlaylistContributors{UserID: userID}); err != nil {
		return err
	}

	return nil
}

// RemoveContributor removes a contributor from a playlist
func (m *PlaylistDBModel) RemoveContributor(playlistID, userID uint) error {
	var playlist Playlist
	if err := m.DB.Preload("Contributors").First(&playlist, playlistID).Error; err != nil {
		return err
	}

	// Assuming you have a relationship set up between Playlist and User models for contributors
	// Adjust the code accordingly based on your database schema
	for i, contributor := range playlist.Contributors {
		if contributor.UserID == userID {
			playlist.Contributors = append(playlist.Contributors[:i], playlist.Contributors[i+1:]...)
			break
		}
	}

	if err := m.DB.Save(&playlist).Error; err != nil {
		return err
	}

	return nil
}

// GetContributors retrieves contributors of a playlist
func (m *PlaylistDBModel) GetContributors(playlistID uint) ([]User, error) {
	var playlist Playlist
	if err := m.DB.Preload("Contributors").First(&playlist, playlistID).Error; err != nil {
		return nil, err
	}

	// Assuming you have a relationship set up between Playlist and User models for contributors
	// Adjust the code accordingly based on your database schema
	var contributors []User
	for _, contributor := range playlist.Contributors {
		contributors = append(contributors, User{UserID: contributor.UserID, Username: "Contributor"}) // Replace with actual user details
	}

	return contributors, nil
}

// AddRelatedPlaylist adds a related playlist to the current playlist
func (m *PlaylistDBModel) AddRelatedPlaylist(playlistID, relatedPlaylistID uint) error {
	var playlist Playlist
	if err := m.DB.Preload("RelatedPlaylists").First(&playlist, playlistID).Error; err != nil {
		return err
	}

	// Assuming you have a relationship set up between Playlist and RelatedPlaylist models
	// Adjust the code accordingly based on your database schema
	if err := m.DB.Model(&playlist).Association("RelatedPlaylists").Append(&RelatedPlaylist{RelatedPlaylistID: relatedPlaylistID}); err != nil {
		return err
	}

	return nil
}

// RemoveRelatedPlaylist removes a related playlist from the current playlist
func (m *PlaylistDBModel) RemoveRelatedPlaylist(playlistID, relatedPlaylistID uint) error {
	var playlist Playlist
	if err := m.DB.Preload("RelatedPlaylists").First(&playlist, playlistID).Error; err != nil {
		return err
	}

	// Assuming you have a relationship set up between Playlist and RelatedPlaylist models
	// Adjust the code accordingly based on your database schema
	for i, relatedPlaylist := range playlist.RelatedPlaylists {
		if relatedPlaylist.RelatedPlaylistID == relatedPlaylistID {
			playlist.RelatedPlaylists = append(playlist.RelatedPlaylists[:i], playlist.RelatedPlaylists[i+1:]...)
			break
		}
	}

	if err := m.DB.Save(&playlist).Error; err != nil {
		return err
	}

	return nil
}

// UpdatePlaylistDetails updates the details of a playlist
func (m *PlaylistDBModel) UpdatePlaylistDetails(playlistID uint, updatedDetails Playlist) error {
	var playlist Playlist
	if err := m.DB.First(&playlist, playlistID).Error; err != nil {
		return err
	}

	// Update playlist details
	playlist.Name = updatedDetails.Name
	playlist.Description = updatedDetails.Description
	playlist.Title = updatedDetails.Title
	playlist.Tags = updatedDetails.Tags
	playlist.IsPublic = updatedDetails.IsPublic
	playlist.FeaturedArtwork = updatedDetails.FeaturedArtwork
	playlist.IsPlayable = updatedDetails.IsPlayable
	playlist.PrivacySetting = updatedDetails.PrivacySetting

	// Save the updated playlist
	if err := m.DB.Save(&playlist).Error; err != nil {
		return err
	}

	return nil
}

// DeletePlaylist deletes a playlist and its associated records
func (m *PlaylistDBModel) DeletePlaylist2(playlistID uint) error {
	var playlist Playlist
	if err := m.DB.First(&playlist, playlistID).Error; err != nil {
		return err
	}

	// Assuming you have cascading delete set up for related records (e.g., Videos, Advertisements)
	// Adjust the code accordingly based on your database schema
	if err := m.DB.Delete(&playlist).Error; err != nil {
		return err
	}

	return nil
}

func (m *PlaylistDBModel) DeletePlaylist3(playlistID uint) error {
	return m.DB.Delete(&Playlist{}, playlistID).Error
}

// GetPlaylistByID retrieves a playlist by ID
func (m *PlaylistDBModel) GetPlaylistByID3(playlistID uint) (*Playlist, error) {
	var playlist Playlist
	if err := m.DB.Preload("Videos").Preload("Advertisements").First(&playlist, playlistID).Error; err != nil {
		return nil, err
	}

	return &playlist, nil
}

// GetPlaybackStatus retrieves the current playback status.
func (p *PlaylistDBModel) GetPlaybackStatus(playlistID uint) string {
	// Add logic to retrieve the current playback status
	// Placeholder code to demonstrate the idea:
	status := "Status: Playing"
	return status
	//return "Status: Playing" // Replace this with actual logic
}

// GetPlaylistInfo retrieves information about the current playlist.
func (p *PlaylistDBModel) GetPlaylistInfo(playlistID uint) string {
	// Add logic to retrieve information about the current playlist
	// Placeholder code to demonstrate the idea:
	playlistInfo := "Playlist: My Playlist"
	return playlistInfo
	//return "Playlist: My Playlist" // Replace this with actual logic
}

// GetVideoQueue retrieves the queue of upcoming videos in the playlist.
func (m *PlaylistDBModel) GetVideoQueue(playlistID uint) ([]Video, error) {
	// Implement logic to get the queue of upcoming videos in the playlist
	// Placeholder code to demonstrate the idea:
	// This could involve returning the remaining videos in the playlist.
	playlist, err := m.GetPlaylistByID(playlistID)
	if err != nil || len(playlist.Videos) == 0 {
		return nil, errors.New("no videos in the playlist")
	}

	// Get the remaining videos in the playlist
	remainingVideos := playlist.Videos[playlist.Order+1:]

	return remainingVideos, nil
}

// IsVideoPlaying checks if a video is currently playing in the playlist.
func (m *PlaylistDBModel) IsVideoPlaying(playlistID uint) bool {
	// Implement logic to check if a video is currently playing in the playlist
	// Placeholder code to demonstrate the idea:
	// This could involve checking the IsPlaying field of the playlist.
	playlist, err := m.GetPlaylistByID(playlistID)
	if err != nil {
		return false
	}
	return playlist.IsPlaying
}

// AdjustVolume adjusts the volume based on the delta value.
func (m *PlaylistDBModel) AdjustVolume(playlistID uint, delta int) int {
	// Implement logic to adjust the volume
	// Placeholder code to demonstrate the idea:
	// This could involve updating the volume level in the playlist.
	playlist, err := m.GetPlaylistByID(playlistID)
	if err != nil {
		return 0
	}

	// Adjust the volume (The actual implementation may vary based on your data model)
	// For simplicity, we'll just increment the volume level based on the delta.
	playlist.Order += delta

	// Save the changes to the database
	if err := m.DB.Save(playlist).Error; err != nil {
		return 0
	}

	return playlist.Order
}

// SkipToPosition skips to the specified position in the playlist.
func (m *PlaylistDBModel) SkipToPosition(playlistID uint, position int) (int, error) {
	// Implement logic to skip to a specific position in the playlist
	// Placeholder code to demonstrate the idea:
	// This could involve updating the order or index of videos in the playlist
	// and returning the new position.

	// Retrieve the playlist
	playlist, err := m.GetPlaylistByID(playlistID)
	if err != nil {
		return 0, err
	}

	// Ensure the position is within the valid range
	if position < 0 || position >= len(playlist.Videos) {
		return 0, errors.New("invalid position")
	}

	// Update the order or index of videos in the playlist
	// (The actual implementation may vary based on your data model)
	// For simplicity, we'll just set the order based on the position for this example.
	for i := range playlist.Videos {
		playlist.Videos[i].Order = position + i
	}

	// Save the changes to the database
	if err := m.DB.Save(playlist).Error; err != nil {
		return 0, err
	}

	return position, nil
}

// AddVideoToPlaylist adds a video to the playlist.
func (m *PlaylistDBModel) AddVideoToPlaylist(playlistID uint, video *Video) error {
	// Implement logic to add a video to the playlist
	// Placeholder code to demonstrate the idea:
	// This could involve updating the Videos field of the playlist.
	playlist, err := m.GetPlaylistByID(playlistID)
	if err != nil {
		return err
	}

	// Ensure the video is not already in the playlist
	for _, existingVideo := range playlist.Videos {
		if existingVideo.ID == video.ID {
			return errors.New("video already in the playlist")
		}
	}

	// Append the video to the playlist
	playlist.Videos = append(playlist.Videos, *video)

	// Save the changes to the database
	if err := m.DB.Save(playlist).Error; err != nil {
		return err
	}

	return nil
}

// RemoveVideoFromPlaylist removes a video from the playlist.
func (m *PlaylistDBModel) RemoveVideoFromPlaylist(playlistID uint, videoID uint) error {
	// Implement logic to remove a video from the playlist
	// Placeholder code to demonstrate the idea:
	// This could involve updating the Videos field of the playlist.
	playlist, err := m.GetPlaylistByID(playlistID)
	if err != nil {
		return err
	}

	// Find the index of the video in the playlist
	var index int
	for i, existingVideo := range playlist.Videos {
		if existingVideo.ID == videoID {
			index = i
			break
		}
	}

	// Remove the video from the playlist
	playlist.Videos = append(playlist.Videos[:index], playlist.Videos[index+1:]...)

	// Save the changes to the database
	if err := m.DB.Save(playlist).Error; err != nil {
		return err
	}

	return nil
}

// ShufflePlaylist shuffles the order of the playlist.
func (m *PlaylistDBModel) ShufflePlaylist(playlistID uint) error {
	// Implement logic to shuffle the playlist order
	// Placeholder code to demonstrate the idea:
	// This could involve randomizing the order of videos in the playlist.
	playlist, err := m.GetPlaylistByID(playlistID)
	if err != nil {
		return err
	}

	// Shuffle the order of videos using Fisher-Yates algorithm
	rand.Seed(time.Now().UnixNano())
	n := len(playlist.Videos)
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		playlist.Videos[i], playlist.Videos[j] = playlist.Videos[j], playlist.Videos[i]
	}

	// Save the changes to the database
	if err := m.DB.Save(playlist).Error; err != nil {
		return err
	}

	return nil
}

// GetCurrentlyPlayingVideoID retrieves the ID of the currently playing video.
func (m *PlaylistDBModel) GetCurrentlyPlayingVideoID(playlistID uint) uint {
	// Implement logic to get the ID of the currently playing video
	// Placeholder code to demonstrate the idea:
	// This could involve checking the currently playing video in the playlist.
	playlist, err := m.GetPlaylistByID(playlistID)
	if err != nil || !playlist.IsPlaying || len(playlist.Videos) == 0 {
		return 0
	}

	return playlist.Videos[playlist.Order].ID
}

// GetCurrentVideoInfo retrieves information about the currently playing video.
func (p *PlaylistDBModel) GetCurrentVideoInfo(playlistID uint) string {
	// Add logic to retrieve information about the currently playing video
	// Placeholder code to demonstrate the idea:
	currentVideoInfo := "Currently playing: Video 1"
	return currentVideoInfo
	//return "Currently playing: Video 1" // Replace this with actual logic
}

// PlayVideo plays a specific video in the playlist.
func (p *PlaylistDBModel) PlayVideo(playlistID uint, videoID uint) (string, error) {
	// Implement logic to play a specific video in the playlist
	// Example: Play the specified video
	videoURL, err := p.VideoDBModel.GetVideoURL(videoID)
	if err != nil {
		return "", fmt.Errorf("Error getting video URL: %v", err)
	}

	// Placeholder code to demonstrate the idea:
	p.VideoPlayer.Play(videoURL)
	return fmt.Sprintf("Now playing %v", videoID), nil
}

// PauseVideo pauses the currently playing video.
func (p *PlaylistDBModel) PauseVideo(playlistID uint) {
	// Implement logic to pause the currently playing video
	// Example: Pause the currently playing video
	p.VideoPlayer.Pause()
}

// ResumeVideo resumes playback of the paused video.
func (p *PlaylistDBModel) ResumeVideo(playlistID uint) {
	// Implement logic to resume playback of the paused video
	// Example: Resume playback of the paused video
	p.VideoPlayer.Resume()
}
