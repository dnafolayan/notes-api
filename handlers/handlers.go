package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Note struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func openDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./notes.db")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func respondWithErr(context *gin.Context, statusCode int, err error) {
	context.JSON(statusCode, gin.H{
		"error": err.Error(),
	})
}

func respondWithCustomErr(context *gin.Context, statusCode int, msg string) {
	context.JSON(statusCode, gin.H{
		"error": msg,
	})
}

func GetNotes(context *gin.Context) {
	db := openDB()
	defer db.Close()

	query := `SELECT id, description, completed FROM notes`

	rows, err := db.Query(query)
	if err != nil {
		respondWithErr(context, http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	var notes []Note

	for rows.Next() {
		var note Note
		var completedInt int

		if err := rows.Scan(&note.ID, &note.Description, &completedInt); err != nil {
			respondWithErr(context, http.StatusInternalServerError, err)
			return
		}

		note.Completed = completedInt == 1
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		respondWithErr(context, http.StatusInternalServerError, err)
		return
	}

	context.JSON(http.StatusOK, notes)
}

func GetNoteByID(context *gin.Context) {
	ID := context.Param("id")

	db := openDB()
	defer db.Close()

	query := `SELECT id, description, completed FROM notes WHERE id = ?`

	var note Note
	var completedInt int

	if err := db.QueryRow(query, ID).Scan(&note.ID, &note.Description, &completedInt); err != nil {
		if err == sql.ErrNoRows {
			respondWithCustomErr(context, http.StatusNotFound, "note not found")
			return
		}

		respondWithErr(context, http.StatusInternalServerError, err)
		return
	}

	note.Completed = completedInt == 1

	context.JSON(http.StatusOK, gin.H{
		"message": "successful",
		"note":    note,
	})
}

func CreateNote(context *gin.Context) {
	db := openDB()

	defer db.Close()

	body, err := io.ReadAll(context.Request.Body)

	if err != nil || len(body) == 0 {
		respondWithCustomErr(context, http.StatusBadRequest, "body cannot be empty")
		return
	}

	var note Note

	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&note); err != nil {
		respondWithCustomErr(context, http.StatusBadRequest, "missing field")
		return
	}
	/*
	 if err := context.BindJSON(&newNote); err != nil {
	 	context.JSON(http.StatusBadRequest, gin.H{
	 		"error": "something went wrong",
	 	})
	 	return
	 }
	*/

	if note.Description == "" {
		respondWithCustomErr(context, http.StatusBadRequest, "missing field")
		return
	}

	query := `INSERT INTO notes (description, completed) VALUES (?, ?)`

	if _, err = db.Exec(query, note.Description, note.Completed); err != nil {
		respondWithErr(context, http.StatusInternalServerError, err)
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Note added successfully",
		"note":    note, // go takes care of pointer dereferencing
	})
}

func ToggleCompleted(context *gin.Context) {
	ID := context.Param("id")

	db := openDB()
	defer db.Close()

	query := `SELECT completed FROM notes WHERE id = ?`

	var status int
	if err := db.QueryRow(query, ID).Scan(&status); err != nil {
		if err == sql.ErrNoRows {
			respondWithCustomErr(context, http.StatusNotFound, "note not found")
			return
		}

		respondWithErr(context, http.StatusInternalServerError, err)
		return
	}

	newStatus := 1
	if status == 1 {
		newStatus = 0
	}

	if _, err := db.Exec(`UPDATE notes SET completed = ? WHERE id = ?`, newStatus, ID); err != nil {
		respondWithErr(context, http.StatusInternalServerError, err)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message":   "successful",
		"completed": newStatus == 1,
	})
}

func UpdateDescription(context *gin.Context) {
	ID := context.Param("id")

	db := openDB()
	defer db.Close()

	var description string

	query := `SELECT description FROM notes WHERE id = ?`
	if err := db.QueryRow(query, ID).Scan(&description); err != nil {
		if err == sql.ErrNoRows {
			respondWithCustomErr(context, http.StatusNotFound, "note not found")
			return
		}

		respondWithErr(context, http.StatusInternalServerError, err)
	}

	type DescriptionInput struct {
		Description string `json:"description" binding:"required"`
	}

	var descriptionInput DescriptionInput

	if err := context.BindJSON(&descriptionInput); err != nil {
		respondWithCustomErr(context, http.StatusBadRequest, "missing field")
		return
	}

	if _, err := db.Exec(`UPDATE notes SET description = ? WHERE id = ?`, descriptionInput.Description, ID); err != nil {
		respondWithErr(context, http.StatusInternalServerError, err)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message":   "successful",
		"completed": descriptionInput,
	})
}

func DeleteNote(context *gin.Context) {
	ID := context.Param("id")

	db := openDB()
	defer db.Close()

	query := `SELECT id FROM notes WHERE id = ?`

	var noteID int

	if err := db.QueryRow(query, ID).Scan(&noteID); err != nil {
		if err == sql.ErrNoRows {
			respondWithCustomErr(context, http.StatusNotFound, "note not found")
			return
		}

		respondWithErr(context, http.StatusInternalServerError, err)
		return
	}

	if _, err := db.Exec(`DELETE FROM notes WHERE id = ?;`, ID); err != nil {
		respondWithErr(context, http.StatusInternalServerError, err)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "deleted successfully",
		"noteID":  noteID,
	})
}
