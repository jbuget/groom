package handlers

import (
	"database/sql"
	"groom/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/meet/v2"
)

// Handler pour lister les rooms en JSON
func ListRoomsJSONHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rooms, err := models.GetAllRooms(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve rooms"})
			return
		}
		c.JSON(http.StatusOK, rooms)
	}
}

// Handler pour créer une room
func CreateRoomHandler(db *sql.DB, meetService *meet.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody struct {
			Slug string `json:"slug"`
		}
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		existingRoom, err := models.GetRoomBySlug(db, requestBody.Slug)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error verifying room existence"})
			return
		}
		if existingRoom != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "A room with the same slug already exists"})
			return
		}

		space, err := meetService.Spaces.Create(&meet.Space{}).Do()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating Google Meet space"})
			return
		}

		room := models.Room{
			Slug:    requestBody.Slug,
			SpaceID: space.MeetingCode,
		}

		id, err := models.CreateRoom(db, room)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting room"})
			return
		}

		room.ID = id
		c.JSON(http.StatusCreated, room)
	}
}

// Handler pour mettre à jour une room
func UpdateRoomHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		var room models.Room
		if err := c.BindJSON(&room); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Mettre à jour la room dans la base de données
		err = models.UpdateRoom(db, id, room)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating room"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Room updated successfully"})
	}
}

// Handler pour supprimer une room
func DeleteRoomHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		// Supprimer la room de la base de données
		err = models.DeleteRoom(db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting room"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Room deleted successfully"})
	}
}
