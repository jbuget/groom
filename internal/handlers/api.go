package handlers

import (
	"database/sql"
	"groom/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
func CreateRoomHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var room models.Room
		if err := c.BindJSON(&room); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Créer la room dans la base de données
		id, err := models.CreateRoom(db, room)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting room"})
			return
		}

		// Retourner la room créée avec son ID
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
