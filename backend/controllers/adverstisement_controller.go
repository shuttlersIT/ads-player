// backend/controllers/advertisement_controller.go

package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/robfig/cron/v3"
	"github.com/shuttlersit/ads-player/backend/models"
)

// AdvertisementController handles advertisement scheduling and playback
type AdvertisementController struct {
	PlaylistModel      *models.PlaylistModel
	AdvertisementModel *models.AdvertisementDBModel
	PlaybackService    PlaybackService
}

// NewAdvertisementController creates a new instance of AdvertisementController
func NewAdvertisementController(playlistModel *models.PlaylistModel, advertisementModel *models.AdvertisementDBModel, playbackService PlaybackService) *AdvertisementController {
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

// UpdateAndRefreshAdvertisements updates and refreshes advertisements from an external source
func (ac *AdvertisementController) UpdateAndRefreshAdvertisements() error {
	// Fetch new advertisements from an external source
	newAdvertisements, err := ac.fetchNewAdvertisements()
	if err != nil {
		log.Printf("Error fetching new advertisements: %v\n", err)
		return err
	}

	// Get the existing advertisements from your storage (e.g., database)
	existingAdvertisements, err := ac.AdvertisementModel.GetAllAdvertisements()
	if err != nil {
		log.Printf("Error fetching existing advertisements: %v\n", err)
		return err
	}

	// Map existing advertisements by their ID for efficient update comparison
	existingAdvertisementMap := make(map[uint]*models.Advertisement)
	for _, existingAd := range existingAdvertisements {
		existingAdvertisementMap[existingAd.ID] = &existingAd
	}

	// Update existing advertisements and add new ones
	for _, newAd := range newAdvertisements {
		if existingAd, exists := existingAdvertisementMap[newAd.ID]; exists {
			// Advertisement already exists, check for updates
			if !areAdvertisementsEqual(existingAd, &newAd) {
				// Advertisement has changed, update it
				err := ac.AdvertisementModel.UpdateAdvertisement(&newAd)
				if err != nil {
					log.Printf("Error updating existing advertisement (ID: %d): %v\n", newAd.ID, err)
				}
			}
		} else {
			// Advertisement is new, create it
			err := ac.AdvertisementModel.CreateAdvertisement(&newAd)
			if err != nil {
				log.Printf("Error creating new advertisement (ID: %d): %v\n", newAd.ID, err)
			}
		}
	}

	// Optionally, you can also delete any advertisements that are not present in the new set
	// (depending on your use case and business logic)

	// Log the refresh event
	log.Printf("Advertisements refreshed at: %v\n", time.Now())

	return nil
}

// fetchNewAdvertisements fetches new advertisements from an AWS S3 bucket
func (ac *AdvertisementController) fetchNewAdvertisements() ([]models.Advertisement, error) {
	// Set your AWS credentials and region
	awsRegion := "your-aws-region"
	awsAccessKeyID := "your-access-key-id"
	awsSecretAccessKey := "your-secret-access-key"

	// Set your S3 bucket name and key (path to the JSON file containing advertisements)
	s3BucketName := "your-s3-bucket-name"
	s3ObjectKey := "path/to/advertisements.json"

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	})
	if err != nil {
		log.Printf("Error creating AWS session: %v\n", err)
		return nil, err
	}

	// Create an S3 client
	s3Client := s3.New(sess)

	// Fetch the JSON file from the S3 bucket
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(s3ObjectKey),
	}

	result, err := s3Client.GetObject(getObjectInput)
	if err != nil {
		log.Printf("Error fetching JSON file from S3: %v\n", err)
		return nil, err
	}
	defer result.Body.Close()

	// Decode the JSON response into a slice of Advertisement structs
	var newAdvertisements []models.Advertisement
	err = json.NewDecoder(result.Body).Decode(&newAdvertisements)
	if err != nil {
		log.Printf("Error decoding JSON response: %v\n", err)
		return nil, err
	}

	return newAdvertisements, nil
}

// fetchNewAdvertisements simulates fetching new advertisements from an external source
func (ac *AdvertisementController) fetchNewAdvertisements2() ([]models.Advertisement, error) {
	// Simulate fetching new advertisements from an external source (e.g., API)
	// Replace the URL with the actual endpoint to fetch advertisements from
	apiURL := "https://example.com/api/advertisements"

	// Make a GET request to the API
	response, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch advertisements, status code: %d", response.StatusCode)
	}

	// Decode the JSON response into a slice of Advertisement structs
	var newAdvertisements []models.Advertisement
	err = json.NewDecoder(response.Body).Decode(&newAdvertisements)
	if err != nil {
		return nil, err
	}

	return newAdvertisements, nil
}

// updateExistingAdvertisements updates existing advertisements based on new data
func (ac *AdvertisementController) updateExistingAdvertisements(newAdvertisements []models.Advertisement) error {
	// Get all existing advertisements from the database
	existingAdvertisements, err := ac.AdvertisementModel.GetAllAdvertisements()
	if err != nil {
		return err
	}

	// Compare existing and new advertisements and update as necessary
	for _, newAd := range newAdvertisements {
		// Check if the new advertisement already exists in the database
		existingAd, exists := ac.findExistingAdvertisement(existingAdvertisements, newAd.Title)
		if exists {
			// Update existing advertisement with new data
			existingAd.Duration = newAd.Duration
			// Update other fields as needed
			// ...

			// Save the updated advertisement to the database
			err := ac.AdvertisementModel.UpdateAdvertisement(existingAd)
			if err != nil {
				return err
			}
			fmt.Printf("Updated existing advertisement: %s\n", existingAd.Title)
		} else {
			// If the advertisement doesn't exist, add it to the database
			err := ac.AdvertisementModel.CreateAdvertisement(&newAd)
			if err != nil {
				return err
			}
			fmt.Printf("Added new advertisement: %s\n", newAd.Title)
		}
	}

	// Remove any outdated advertisements not present in the new data
	err = ac.removeOutdatedAdvertisements(existingAdvertisements, newAdvertisements)
	if err != nil {
		return err
	}

	return nil
}

// findExistingAdvertisement finds an existing advertisement by name in a slice of advertisements
func (ac *AdvertisementController) findExistingAdvertisement(advertisements []models.Advertisement, name string) (*models.Advertisement, bool) {
	for _, ad := range advertisements {
		if ad.Title == name {
			return &ad, true
		}
	}
	return nil, false
}

// removeOutdatedAdvertisements removes outdated advertisements not present in the new data
func (ac *AdvertisementController) removeOutdatedAdvertisements(existingAdvertisements []models.Advertisement, newAdvertisements []models.Advertisement) error {
	// Identify outdated advertisements by comparing existing and new data
	outdatedAdvertisements := ac.findOutdatedAdvertisements(existingAdvertisements, newAdvertisements)

	// Remove outdated advertisements from the database
	for _, outdatedAd := range outdatedAdvertisements {
		err := ac.AdvertisementModel.DeleteAdvertisement(outdatedAd.ID)
		if err != nil {
			return err
		}
		fmt.Printf("Removed outdated advertisement: %s\n", outdatedAd.Title)
	}

	return nil
}

// findOutdatedAdvertisements identifies outdated advertisements not present in the new data
func (ac *AdvertisementController) findOutdatedAdvertisements(existingAdvertisements []models.Advertisement, newAdvertisements []models.Advertisement) []models.Advertisement {
	var outdatedAdvertisements []models.Advertisement

	// Identify outdated advertisements by comparing existing and new data
	for _, existingAd := range existingAdvertisements {
		found := false
		for _, newAd := range newAdvertisements {
			if existingAd.Title == newAd.Title {
				found = true
				break
			}
		}
		if !found {
			outdatedAdvertisements = append(outdatedAdvertisements, existingAd)
		}
	}

	return outdatedAdvertisements
}

// areAdvertisementsEqual checks if two advertisements are equal
func areAdvertisementsEqual(ad1, ad2 *models.Advertisement) bool {
	// Compare all relevant fields to determine equality
	return ad1.ID == ad2.ID &&
		ad1.Name == ad2.Name &&
		ad1.Duration == ad2.Duration &&
		ad1.ContentURL == ad2.ContentURL &&
		ad1.Description == ad2.Description &&
		ad1.StartDate.Equal(ad2.StartDate) && // Compare StartDate
		ad1.EndDate.Equal(ad2.EndDate) && // Compare EndDate
		ad1.TargetAudience == ad2.TargetAudience &&
		// Add more field comparisons based on your Advertisement struct
		areExternalLinksEqual(ad1.ExternalLinks, ad2.ExternalLinks) &&
		areMediaAttachmentsEqual(ad1.MediaAttachments, ad2.MediaAttachments) &&
		areAdvertisementHashtagsEqual(ad1.AdvertisementHashtags, ad2.AdvertisementHashtags) &&
		areAdvertisementAnalyticsEqual(ad1.Analytics, ad2.Analytics)
	// Add more field comparisons based on your Advertisement struct
}

// areExternalLinksEqual checks if two slices of ExternalLink structs are equal
func areExternalLinksEqual(links1, links2 []models.ExternalLink) bool {
	if len(links1) != len(links2) {
		return false
	}

	// Compare each ExternalLink
	for i := range links1 {
		if links1[i].URL != links2[i].URL || links1[i].Label != links2[i].Label {
			return false
		}
	}

	return true
}

// areMediaAttachmentsEqual checks if two slices of MediaAttachment structs are equal
func areMediaAttachmentsEqual(attachments1, attachments2 []models.MediaAttachment) bool {
	if len(attachments1) != len(attachments2) {
		return false
	}

	// Compare each MediaAttachment
	for i := range attachments1 {
		if attachments1[i].Type != attachments2[i].Type || attachments1[i].URL != attachments2[i].URL {
			return false
		}
	}

	return true
}

// areHashtagsEqual checks if two slices of Hashtag structs are equal
func areAdvertisementHashtagsEqual(tags1, tags2 []models.AdvertisementHashtag) bool {
	if len(tags1) != len(tags2) {
		return false
	}

	// Compare each Hashtag
	for i := range tags1 {
		if tags1[i].Tag != tags2[i].Tag {
			return false
		}
	}

	return true
}

// areHashtagsEqual checks if two slices of Hashtag structs are equal
func areHashtagsEqual(tags1, tags2 []models.Hashtag) bool {
	if len(tags1) != len(tags2) {
		return false
	}

	// Compare each Hashtag
	for i := range tags1 {
		if tags1[i].Tag != tags2[i].Tag || tags1[i].Tags != tags2[i].Tags { // Compare Tags field
			return false
		}
	}

	return true
}

// areAdvertisementAnalyticsEqual checks if two AdvertisementAnalytics structs are equal
func areAdvertisementAnalyticsEqual(analytics1, analytics2 models.AdvertisementAnalytics) bool {
	return analytics1.Views == analytics2.Views &&
		analytics1.Clicks == analytics2.Clicks &&
		analytics1.Conversions == analytics2.Conversions && // Added Conversions field
		analytics1.Revenue == analytics2.Revenue // Added Revenue field
	// Add more field comparisons based on your AdvertisementAnalytics struct
}
