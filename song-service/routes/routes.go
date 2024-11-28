package routes

import (
	"song-service/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	songs := r.Group("/api/v1")
	{
		songs.GET("/songs/search/filters", controllers.GetSongs)
		songs.GET("/songs/search/verse", controllers.GetVerse)
		songs.POST("/songs", controllers.AddSong)
		songs.PUT("/songs/:id", controllers.UpdateSong)
		songs.DELETE("/songs/:id", controllers.DeleteSong)
	}
}
