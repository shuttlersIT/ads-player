// backend/routes/video_routes.go

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shuttlersit/ads-player/backend/controllers"
	"github.com/shuttlersit/ads-player/backend/services"
)

// SetVideoRoutes sets up video-related routes.
func SetVideoRoutes(r *gin.Engine, videoService *services.VideoService) {
	videoController := controllers.NewVideoController(videoService)

	// ... (existing routes)

	// New routes
	r.GET("/videos/:id/related", videoController.GetRelatedVideos)
	r.GET("/videos/most-liked", videoController.GetMostLikedVideos)
	r.GET("/users/:id/videos", videoController.GetVideosByUploader)
	r.GET("/categories/:id/videos", videoController.GetVideosByCategory)
	r.PATCH("/videos/:id/increment-views", videoController.IncrementViews)
	r.PATCH("/videos/:id/increment-likes", videoController.IncrementLikes)
	r.PATCH("/videos/:id/increment-dislikes", videoController.IncrementDislikes)
	r.GET("/videos/top-tags", videoController.GetTopTags)
	// New routes
	r.POST("/videos", videoController.CreateVideo)
	r.PUT("/videos/:id", videoController.UpdateVideo)
	r.GET("/videos/:id/analytics", videoController.GetVideoAnalytics)
}
