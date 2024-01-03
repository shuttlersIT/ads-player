// backend/controllers/remote_control_controller.go

package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/shuttlersit/ads-player/backend/models"
	"github.com/shuttlersit/ads-player/backend/services"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// RemoteControlController handles remote control via WebSocket
type RemoteControlController struct {
	PlaylistService *services.DefaultPlaylistService // Add PlaylistService dependency
	PlaybackService *services.PlaybackService        // Add PlaybackService dependency
	PlaylistModel   *models.PlaylistDBModel          // Add PlaylistService dependency
}

// NewRemoteControlController creates a new instance of RemoteControlController
func NewRemoteControlController(playlistService *services.DefaultPlaylistService, playbackService *services.PlaybackService, PlaylistModel *models.PlaylistDBModel) *RemoteControlController {
	return &RemoteControlController{
		PlaylistService: playlistService,
		PlaybackService: playbackService,
		PlaylistModel:   PlaylistModel,
	}
}

// APIControlHandler handles remote control via API
func (rc *RemoteControlController) APIControlHandler(c *gin.Context) {
	command := c.Query("command")

	switch command {
	case "next":
		rc.PlaylistService.PlayNextVideo() // Add logic to play the next video
		c.JSON(http.StatusOK, gin.H{"result": "Playing the next video"})
	case "pause":
		rc.PlaylistService.PausePlayback() // Add logic to pause playback
		c.JSON(http.StatusOK, gin.H{"result": "Playback paused"})
	// Add more command cases as needed

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown command"})
	}
}

// WebSocketHandler handles WebSocket connections for remote control
func (rc *RemoteControlController) WebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Handle WebSocket messages
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading WebSocket message:", err)
			break
		}

		// Handle the received message
		command := string(p)
		response := rc.handleWebSocketCommand(command)

		// You can send a response if needed
		if err := conn.WriteMessage(messageType, []byte(response)); err != nil {
			fmt.Println("Error writing WebSocket message:", err)
			break
		}
	}
}

// handleWebSocketCommand handles the WebSocket command and returns a response
func (rc *RemoteControlController) handleWebSocketCommand(command string) string {
	switch command {
	case "next":
		rc.PlaybackService.PlayNextVideo()
		return "Playing the next video"
	case "pause":
		rc.PlaybackService.PausePlayback()
		return "Playback paused"
	case "get_status":
		status := rc.PlaybackService.GetPlaybackStatus()
		return status
	case "get_playlist":
		playlistInfo := rc.PlaylistService.GetPlaylistInfo()
		return playlistInfo
	case "volume_up":
		rc.PlaylistService.AdjustVolume(1)
		return "Volume increased"
	case "volume_down":
		rc.PlaylistService.AdjustVolume(-1)
		return "Volume decreased"
	case "skip_to":
		position := parsePositionFromCommand(command)
		rc.PlaylistService.SkipToPosition(position)
		return "Skipped to the specified position"
	case "add_to_playlist":
		videoID := parseVideoIDFromCommand(command)
		rc.PlaylistService.AddToPlaylist(videoID)
		return "Added to the playlist"
	case "remove_from_playlist":
		videoID := parseVideoIDFromCommand(command)
		rc.PlaylistService.RemoveFromPlaylist(videoID)
		return "Removed from the playlist"
	case "shuffle_playlist":
		rc.PlaylistService.ShufflePlaylist()
		return "Playlist shuffled"
	case "get_current_video":
		currentVideoInfo := rc.PlaybackService.GetCurrentVideoInfo()
		return currentVideoInfo
	case "play_video":
		videoID := parseVideoIDFromCommand(command)
		rc.PlaylistService.PlayVideo(videoID)
		return "Playing the specified video"
	case "pause_video":
		rc.PlaylistService.PauseVideo()
		return "Video paused"
	case "resume_video":
		rc.PlaylistService.ResumeVideo()
		return "Video resumed"
	default:
		return "Unknown command"
	}
}

// handleWebSocketCommand handles the WebSocket command and returns a response
func (rc *RemoteControlController) handleWebSocketCommand2(command string) string {
	switch command {
	case "next":
		rc.PlaylistService.PlayNextVideo() // Add logic to play the next video
		return "Playing the next video"
	case "pause":
		rc.PlaylistService.PausePlayback() // Add logic to pause playback
		return "Playback paused"
	case "get_status":
		// Add logic to get the current playback status and send it as a response
		status := rc.PlaylistService.GetPlaybackStatus()
		return status
		// return "Status: Playing"
	case "get_playlist":
		// Add logic to get information about the current playlist and send it as a response
		playlistInfo := rc.PlaylistService.GetPlaylistInfo()
		return playlistInfo
		//return "Playlist: My Playlist"
	case "volume_up":
		// Add logic to increase the volume
		rc.PlaylistService.AdjustVolume(1)
		return "Volume increased"
	case "volume_down":
		// Add logic to decrease the volume
		rc.PlaylistService.AdjustVolume(-1)
		return "Volume decreased"
	case "skip_to":
		// Add logic to skip to a specific position in the playlist
		position := parsePositionFromCommand(command)
		playlistID := rc.PlaylistModel.GetCurrentPlaylistID()
		newPosition, err := rc.PlaylistService.SkipToPosition(playlistID, position)
		if err != nil {
			return fmt.Sprintf("Error: %v", err)
		}
		return fmt.Sprintf("Video Position Changed to %d", newPosition)
	case "add_to_playlist":
		// Add logic to add a video or item to the playlist
		videoID := parseVideoIDFromCommand(command)
		rc.PlaylistService.AddToPlaylist(videoID)
		return "Added to the playlist"
		// Add more command cases as needed
	case "remove_from_playlist":
		// Add logic to remove a video or item from the playlist
		videoID := parseVideoIDFromCommand(command)
		rc.PlaylistService.RemoveFromPlaylist(videoID)
		return "Removed from the playlist"
	case "shuffle_playlist":
		// Add logic to shuffle the playlist order
		rc.PlaylistService.ShufflePlaylist()
		return "Playlist shuffled"
	case "get_current_video":
		// Add logic to get information about the currently playing video
		currentVideoInfo := rc.PlaylistService.GetCurrentVideoInfo()
		return currentVideoInfo
		//return "Currently playing: Video 1"
	case "play_video":
		// Add logic to play a specific video in the playlist
		videoID := parseVideoIDFromCommand(command)
		rc.PlaylistService.PlayVideo(videoID)
		return "Playing the specified video"
		//return "Playing the specified video"
	case "pause_video":
		// Add logic to pause the currently playing video
		rc.PlaylistService.PauseVideo()
		return "Video paused"
	case "resume_video":
		// Add logic to resume playback of the paused video
		rc.PlaylistService.ResumeVideo()
		return "Video resumed"

	default:
		return "Unknown command"
	}
}
