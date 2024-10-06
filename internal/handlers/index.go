package handlers

import (
	"database/sql"
	"groom/internal/models"
	"net/http"
	"strings"

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
func ListRoomsHTMLHandler(db *sql.DB, meetService *meet.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		rooms, err := models.GetAllRooms(db)
		if err != nil {
			c.String(http.StatusInternalServerError, "Unable to retrieve rooms")
			return
		}

		activeConferences, err := meetService.ConferenceRecords.List().Filter("end_time IS NULL").Do()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "unhealthy",
				"error":  "Google Meet service unavailable",
			})
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
				IsOccupied: isRoomOccupied(room.SpaceID, activeConferences.ConferenceRecords),
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
func RedirectHandler(db *sql.DB, meetService *meet.Service) gin.HandlerFunc {
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

		var meetingCode string

		// Si le SpaceID commence par "spaces/", on doit récupérer le Space via meetService
		if strings.HasPrefix(room.SpaceID, "spaces/") {
			space, err := meetService.Spaces.Get(room.SpaceID).Do()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Google Meet space", "details": err.Error()})
				return
			}

			// Mettre à jour le SpaceID avec le meeting code (si disponible)
			if space.MeetingCode != "" {
				meetingCode = space.MeetingCode
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Meeting code not found for the space"})
				return
			}
		} else {
			meetingCode = room.SpaceID
		}
		// TODO check GMeet space validity

		// Rediriger vers la room Google Meet correspondante
		c.Redirect(http.StatusFound, "https://meet.google.com/"+meetingCode)
	}
}

// GET /healthz
func HealthzHandler(db *sql.DB, meetService *meet.Service) gin.HandlerFunc {
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
		_, err = meetService.ConferenceRecords.List().Do()
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
