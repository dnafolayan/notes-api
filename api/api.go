package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Note struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

var (
	notes      = []*Note{}
	nextID int = 1
)

func convertIDToString(context *gin.Context, param string) (int, error) {
	idParam := context.Param(param)
	ID, err := strconv.Atoi(idParam)
	if err != nil {
		return 0, err
	}

	return ID, nil
}

func GetNotes(context *gin.Context) {
	context.JSON(http.StatusOK, notes)
}

func GetNoteByID(context *gin.Context) {
	ID, err := convertIDToString(context, "id")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ID",
		})

		return
	}

	for _, note := range notes {
		if note.ID == ID {
			context.JSON(http.StatusOK, gin.H{
				"message": "successful",
				"note":    note,
			})

			return
		}
	}

	context.JSON(http.StatusNotFound, gin.H{
		"error": "note not found",
	})
}

func PostNote(context *gin.Context) {
	body, err := io.ReadAll(context.Request.Body)

	if err != nil || len(body) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "body cannot be empty",
		})

		return
	}

	var newNote Note

	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&newNote); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "missing fields",
		})

		return
	}

	// if err := context.BindJSON(&newNote); err != nil {
	// 	context.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "something went wrong",
	// 	})
	// 	return
	// }

	if newNote.Description == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "missing field",
		})

		return
	}

	newNote.ID = nextID
	nextID++

	notes = append(notes, &newNote)

	context.JSON(http.StatusCreated, gin.H{
		"message": "Note added successfully",
		"note":    newNote, // go takes care of pointer dereferencing
	})
}

func ToggleCompleted(context *gin.Context) {
	ID, err := convertIDToString(context, "id")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ID",
		})

		return
	}

	for i := range notes {
		if notes[i].ID == ID {
			notes[i].Completed = !notes[i].Completed

			context.JSON(http.StatusOK, gin.H{
				"message": "successful",
				"note":    notes[i],
			})

			return
		}
	}

	context.JSON(http.StatusNotFound, gin.H{
		"error": "note not found",
	})
}

func UpdateDescription(context *gin.Context) {
	ID, err := convertIDToString(context, "id")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ID",
		})

		return
	}

	type DescriptionInput struct {
		Description string `json:"description" binding:"required"`
	}

	var descriptionInput DescriptionInput

	if err := context.BindJSON(&descriptionInput); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "missing field",
		})
	}

	for _, note := range notes {
		if note.ID == ID {
			note.Description = descriptionInput.Description

			context.JSON(http.StatusOK, gin.H{
				"message": "successful",
				"note":    note,
			})

			return
		}
	}

	context.JSON(http.StatusNotFound, gin.H{
		"error": "note not found",
	})
}

func DeleteNote(context *gin.Context) {
	ID, err := convertIDToString(context, "id")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ID",
		})

		return
	}

	for i, note := range notes {
		if note.ID == ID {
			notes = slices.Delete(notes, i, i+1)

			context.JSON(http.StatusOK, gin.H{
				"message": "deleted successfully",
				"note":    note,
			})

			return
		}
	}

	context.JSON(http.StatusNotFound, gin.H{
		"error": "note not found",
	})
}
