package routes

import (
	"github.com/gin-gonic/gin"
	"notes-api/controllers"
)

func RegisterNoteRoutes(r *gin.Engine) {
	r.GET("/notes", controllers.GetNotes)
	r.POST("/notes", controllers.CreateNote)
	r.PUT("/notes/:id", controllers.UpdateNote)
}
