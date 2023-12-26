// backend/controllers/remote_control_controller.go

package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
}

// NewRemoteControlController creates a new instance of RemoteControlController
func NewRemoteControlController(playlistService *services.DefaultPlaylistService) *RemoteControlController {
	return &RemoteControlController{
		PlaylistService: playlistService,
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
		rc.PlaylistService.PlayNextVideo() // Add logic to play the next video
		return "Playing the next video"
	case "pause":
		rc.PlaylistService.PausePlayback() // Add logic to pause playback
		return "Playback paused"
	// Add more command cases as needed

	default:
		return "Unknown command"
	}
}
