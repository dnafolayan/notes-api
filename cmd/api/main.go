package main

import (
	"log"

	"github.com/dnafolayan/notes-api/db"
	"github.com/dnafolayan/notes-api/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	router := gin.Default()

	router.GET("/notes", handlers.GetNotes)
	router.GET("/notes/:id", handlers.GetNoteByID)
	router.POST("/notes", handlers.PostNote)
	router.PATCH("/notes/completed/:id", handlers.ToggleCompleted)
	router.PATCH("/notes/description/:id", handlers.UpdateDescription)
	router.DELETE("/notes/delete/:id", handlers.DeleteNote)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
