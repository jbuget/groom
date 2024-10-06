package handlers

import (
	"database/sql"
	googleapi "groom/internal/google"
	"groom/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/meet/v2"
)

// Fonction pour déterminer si une room est occupée
func isRoomOccupied(spaceID string, activeConferences []*meet.ConferenceRecord) bool {
	for _, conference := range activeConferences {
		if conference.Space == spaceID {
			return true
		}
	}
	return false
}

// GET /
func ListRoomsHTMLHandler(db *sql.DB, meetService *googleapi.MeetClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		rooms, err := models.GetAllRooms(db)
		if err != nil {
			c.String(http.StatusInternalServerError, "Unable to retrieve rooms")
			return
		}

		activeConferences, err := meetService.ListActiveConferences()
		if err != nil {
			c.String(http.StatusInternalServerError, "Unable to retrieve active conferences")
			return
		}

		// RoomView est un type dérivé qui contient les informations nécessaires pour l'affichage
		type RoomView struct {
			ID         int    `json:"id"`
			Slug       string `json:"slug"`
			SpaceID    string `json:"space_id"`
			IsOccupied bool   `json:"is_occupied"`
		}

		// Créer une liste de RoomView pour passer au template
		var roomViews []RoomView
		for _, room := range rooms {
			// Créer une instance de RoomView pour chaque Room
			roomView := RoomView{
				ID:      room.ID,
				Slug:    room.Slug,
				SpaceID: room.SpaceID,
				// Ajouter la propriété IsOccupied en fonction des conférences actives
				IsOccupied: isRoomOccupied(room.SpaceID, activeConferences),
			}

			// Ajouter RoomView à la liste
			roomViews = append(roomViews, roomView)
		}

		// Render du template list.html avec la liste des rooms
		c.HTML(http.StatusOK, "list.html", gin.H{
			"rooms": roomViews,
		})
	}
}

// GET /:slug
func RedirectHandler(db *sql.DB, meetService *googleapi.MeetClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")

		room, err := models.GetRoomBySlug(db, slug)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error verifying room existence"})
			return
		}
		if room == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}

		space, err := meetService.GetSpace(room.SpaceID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Google Meet space", "details": err.Error()})
			return
		}

		// Rediriger vers la room Google Meet correspondante
		c.Redirect(http.StatusFound, space.MeetingUri)
	}
}

// GET /healthz
func HealthzHandler(db *sql.DB, meetService *googleapi.MeetClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Vérifier la connexion à la base de données
		err := db.Ping()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "unhealthy",
				"error":  "Database connection failed",
			})
			return
		}

		// Vérifier l'accès à l'API Google Meet
		err = meetService.CheckMeetClient()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "unhealthy",
				"error":  "Google Meet service unavailable",
			})
			return
		}

		// Si tout va bien, renvoyer un statut healthy
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	}
}
