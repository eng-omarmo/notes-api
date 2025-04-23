package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"notes-api/models"
	"strings"
)

func GetNotes(c *gin.Context) {
	var notes []models.Note
	models.DB.Find(&notes)
	c.JSON(http.StatusOK, gin.H{"data": notes})
}

func UpdateNote(c *gin.Context) {
	// 1. Get note ID
	id := c.Param("id")

	// 2. Find the existing note
	var note models.Note
	if err := models.DB.First(&note, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	fmt.Print("Note found:", note)

	// 3. Use separate input struct for binding
	type UpdateNoteInput struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	var input UpdateNoteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 4. Update only provided fields
	if input.Title != "" {
		note.Title = strings.TrimSpace(input.Title)
	}
	if input.Content != "" {
		note.Content = strings.TrimSpace(input.Content)
	}

	// 5. Save changes
	if err := models.DB.Save(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
		return
	}


	c.JSON(http.StatusOK, gin.H{"data": note})
}
func CreateNote(c *gin.Context) {
	// 1. Create dedicated input struct for better validation control
	// type CreateNoteInput struct {
	// 	Title   string `json:"title" binding:"required,min=1"`
	// 	Content string `json:"content" binding:"required,min=1"`
	// }

	type CreateNoteInput struct {
		Title   string `json:"title" binding:"required,min=1"`
		Content string `json:"content" binding:"required,min=1"`
	}

	var input CreateNoteInput

	// 2. Improved error handling with structured messages
	// if err := c.ShouldBindJSON(&input); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error":  "Invalid request payload",
	// 		"detail": err.Error(),
	// 	})
	// 	return
	// }

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Invalid request payload",
			"detail": err.Error(),
		})
	}

	// 3. Check for existing note with proper error handling
	// var existingNote models.Note
	// err := models.DB.
	// 	Where("LOWER(title) = LOWER(?)", input.Title). // Case-insensitive check
	// 	First(&existingNote).
	// 	Error

	var existingNote models.Note
	err := models.DB.
		Where("LOWER(title) = LOWER(?)", input.Title). // Case-insensitive check
		First(&existingNote).
		Error
	if err == nil {
		// 4. Use proper HTTP status code for conflict
		c.JSON(http.StatusConflict, gin.H{
			"error": "Note with this title already exists",
		})

		fmt.Println("Note with this title already exists:", input.Title)
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 5. Handle unexpected database errors
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to verify note uniqueness",
		})
		return
	}

	// 6. Create note with explicit field mapping
	newNote := models.Note{
		Title:   strings.TrimSpace(input.Title),
		Content: strings.TrimSpace(input.Content),
	}

	// 7. Handle database creation errors
	if err := models.DB.Create(&newNote).Error; err != nil {
		// 8. Check for duplicate in case of race condition
		if strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Note with this title was created concurrently",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create note",
		})
		return
	}

	// 9. Use proper HTTP status code for resource creation
	c.JSON(http.StatusCreated, gin.H{
		"data": newNote,
	})
}
