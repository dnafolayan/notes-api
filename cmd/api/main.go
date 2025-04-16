package main

import (
	"log"

	"github.com/dnafolayan/notes_api/api"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/notes", api.GetNotes)
	router.GET("/notes/:id", api.GetNoteByID)
	router.POST("/notes", api.PostNote)
	router.PATCH("/notes/:id", api.ToggleCompleted)
	router.PATCH("/notes/:id", api.ModifyDescription)
	router.DELETE("/notes/:id", api.DeleteNote)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
