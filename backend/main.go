// backend/main.go

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/shuttlersit/ads-player/backend/controllers"
	"github.com/shuttlersit/ads-player/backend/models"
	"github.com/shuttlersit/ads-player/backend/routes"
	"github.com/shuttlersit/ads-player/backend/services"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func main() {

	s3BaseURL := models.S3BaseURL
	bucketURL, region, accessKey, secretKey := "url", "region", "accessKey", "secretKey"
	// Migrate the schema
	db.AutoMigrate(&models.Playlist{}, &models.Video{}, &models.Advertisement{}, &models.AdvertisementPlayEvent{}, &services.DefaultS3Service{})

	// Create models
	logger := services.NewLogger()
	advertisementModel := models.NewAdvertisementDBModel(db)
	videoDBModel := models.NewVideoDBModel(db, s3BaseURL)
	s3Service := services.NewDefaultS3Service(bucketURL, region, accessKey, secretKey)
	videoPlayer := services.NewRTSPVideoPlayer(s3Service) // Initialize the video player
	playlistModel := models.NewPlaylistModel(db, videoDBModel)
	videoService := services.NewVideoService(videoDBModel, videoPlayer)
	database := db
	currentPlaylist := playlistModel.GetCurrentPlaylistID()
	currentAdvertisement := &models.Advertisement{}

	playlistService := services.NewDefaultPlaylistService(playlistModel, videoDBModel, advertisementModel, videoService, currentPlaylist)
	playbackService := services.NewPlaybackService(playlistModel, advertisementModel, s3Service, videoService, playlistService, logger, database, videoDBModel, currentPlaylist, currentAdvertisement)

	// Create controllers
	advertisementController := controllers.NewAdvertisementController(videoDBModel, playlistModel, advertisementModel, playbackService, s3Service)
	playlistController := controllers.NewPlaylistDBController(db, playlistService)
	remoteControlController := controllers.NewRemoteControlController(playlistService, playbackService, playlistModel)
	// Note: PlaylistController object not needed.

	// Start the advertisement scheduler
	schedulerCtx, schedulerCancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := runAdvertisementScheduler(schedulerCtx, advertisementController); err != nil {
			log.Fatal("Error starting Advertisement Scheduler: ", err)
		}
	}()

	// Initialize Gin router
	r := gin.Default()

	// Register playlist routes
	playlistController.RegisterRoutes(r)

	// Register playlist routes
	routes.RegisterPlaylistRoutes(r, db, playlistController, remoteControlController)

	// Register remote control routes
	routes.RegisterRemoteControlRoutes(r, advertisementController, playlistService, playbackService)

	// Set up routes
	routes.SetVideoRoutes(r, videoService)

	// Register video streaming route
	routes.GET("/video/:key", videoPlayer.StreamVideoHandler)

	// Graceful shutdown
	gracefulShutdown(r, schedulerCancel, &wg)

	// Run the server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error running the server: ", err)
	}

	//r.Run(":" + config.PORT)
}

func runAdvertisementScheduler(ctx context.Context, advertisementController *controllers.AdvertisementController) error {
	c := cron.New()

	// Schedule a daily job to update and refresh advertisements
	_, err := c.AddFunc("0 0 * * *", func() {
		err := advertisementController.UpdateAndRefreshAdvertisements()
		if err != nil {
			fmt.Println("Error updating and refreshing advertisements:", err)
		}
	})
	if err != nil {
		return err
	}

	_, err = c.AddFunc("*/5 * * * *", func() {
		// Get playlists that are eligible for advertisements
		playlists, err := advertisementController.PlaylistModel.GetPlaylistsForAdvertisements()
		if err != nil {
			fmt.Println("Error fetching playlists for advertisements:", err)
			return
		}

		// Loop through playlists and schedule advertisements
		for _, playlist := range playlists {
			err := advertisementController.ScheduleAdvertisementForPlaylist(playlist)
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

	// Run the scheduler until the context is canceled
	<-ctx.Done()

	// Stop the cron scheduler gracefully
	c.Stop()

	return nil
}

func gracefulShutdown(router *gin.Engine, cancel context.CancelFunc, wg *sync.WaitGroup) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		fmt.Println("\nShutting down gracefully...")

		// Cancel the context to stop the scheduler
		cancel()

		// Wait for all goroutines to finish
		wg.Wait()

		// Perform any additional cleanup or resource releasing here

		fmt.Println("Graceful shutdown completed.")
		os.Exit(0)
	}()
}
