package handlers

import (
	"database/sql"
	googleapi "groom/internal/google"
	"groom/internal/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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
func CreateRoomHandler(db *sql.DB, meetService *googleapi.MeetClient) gin.HandlerFunc {
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

		space, err := meetService.CreateSpace()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating Google Meet space"})
			return
		}

		room := models.Room{
			Slug:    requestBody.Slug,
			SpaceID: space.Name,
		}

		createdRoom, err := models.CreateRoom(db, room)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting room"})
			return
		}

		c.JSON(http.StatusCreated, createdRoom)
	}
}

func UpdateRoomHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get and check query param ID
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		room, err := models.GetRoomByID(db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying for room"})
			return
		}
		if room == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}

		// Get and check body params
		var requestBody struct {
			Slug    string `json:"slug"`
			SpaceID string `json:"space_id"`
		}
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Update room
		room.Slug = requestBody.Slug
		room.SpaceID = requestBody.SpaceID
		room.UpdatedAt = time.Now()

		err = models.UpdateRoom(db, *room)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating room"})
			return
		}

		// Return
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

func MigrateRoomsHandler(db *sql.DB, meetService *googleapi.MeetClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Récupérer toutes les rooms depuis la base de données
		rooms, err := models.GetAllRooms(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve rooms"})
			return
		}

		// Parcourir chaque room
		for _, room := range rooms {
			// Vérifier si le SpaceID ne commence pas par "spaces/"
			if !strings.HasPrefix(room.SpaceID, "spaces/") {
				// Appeler la fonction GetSpace pour récupérer le Space correspondant
				space, err := meetService.GetSpace("spaces/"+room.SpaceID)
				if err != nil {
					log.Printf("Failed to retrieve space for Room ID %d: %v", room.ID, err)
					continue // Ignorer cette room et continuer avec les autres
				}

				// Mettre à jour la propriété SpaceID avec space.Name
				room.SpaceID = space.Name

				// Mettre à jour la room dans la base de données
				err = models.UpdateRoom(db, room)
				if err != nil {
					// Gestion des erreurs lors de la mise à jour de la room
					log.Printf("Failed to update room ID %d: %v", room.ID, err)
					continue // Ignorer cette room et continuer avec les autres
				}
			}
		}

		// Retourner un succès si toutes les rooms ont été migrées
		c.JSON(http.StatusOK, gin.H{"message": "Rooms migration completed successfully"})
	}
}
