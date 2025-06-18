package handlers

import (
	"database/sql"
	"fmt"
	googleapi "groom/internal/google"
	"groom/internal/models"
	"net/http"
	"sort"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Fonction pour déterminer si une room est occupée
func isRoomOccupied(spaceID string, activeConferences []*googleapi.ConferenceDTO) bool {
	for _, conference := range activeConferences {
		if conference.SpaceID == spaceID {
			return true
		}
	}
	return false
}

func getRoomParticipantCount(spaceID string, activeConferences []*googleapi.ConferenceDTO) int {
	for _, conference := range activeConferences {
		if conference.SpaceID == spaceID {
			return len(conference.Participants)
		}
	}
	return 0
}

// GET /
func ListRoomsHTMLHandler(db *sql.DB, meetService *googleapi.MeetClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from session for star status

		session := sessions.Default(c)
		userID := session.Get("user")
		userId := ""
		if userID != nil {
			userId = userID.(string)
		}

		fmt.Println("session:", session)

		// Get rooms with star status if user is logged in
		var rooms []models.Room
		var err error

		if userId != "" {
			rooms, err = models.GetAllRoomsWithStarStatus(db, userId)
		} else {
			rooms, err = models.GetAllRooms(db)
		}

		if err != nil {
			c.String(http.StatusInternalServerError, "Unable to retrieve rooms")
			return
		}

		activeConferences, err := meetService.ListActiveConferences()
		if err != nil {
			c.String(http.StatusInternalServerError, "Unable to retrieve active conferences")
			return
		}

		type RoomView struct {
			ID               int    `json:"id"`
			Slug             string `json:"slug"`
			SpaceID          string `json:"space_id"`
			IsOccupied       bool   `json:"is_occupied"`
			ParticipantCount int    `json:"participant_count"`
			IsStarred        bool   `json:"is_starred"`
		}

		var allRooms []RoomView

		for _, room := range rooms {
			isOccupied := isRoomOccupied(room.SpaceID, activeConferences)
			roomView := RoomView{
				ID:               room.ID,
				Slug:             room.Slug,
				SpaceID:          room.SpaceID,
				IsOccupied:       isOccupied,
				ParticipantCount: getRoomParticipantCount(room.SpaceID, activeConferences),
				IsStarred:        room.IsStarred,
			}
			allRooms = append(allRooms, roomView)
		}

		// Sort rooms in the new order: Starred+occupied, Starred, Occupied, Rest
		sort.SliceStable(allRooms, func(i, j int) bool {
			roomA, roomB := allRooms[i], allRooms[j]
			
			// Define priority for each room type
			getPriority := func(room RoomView) int {
				if room.IsStarred && room.IsOccupied {
					return 1 // Starred + occupied (highest priority)
				} else if room.IsStarred {
					return 2 // Starred only
				} else if room.IsOccupied {
					return 3 // Occupied only
				} else {
					return 4 // Rest (lowest priority)
				}
			}
			
			priorityA := getPriority(roomA)
			priorityB := getPriority(roomB)
			
			// Sort by priority first
			if priorityA != priorityB {
				return priorityA < priorityB
			}
			
			// Within same priority, sort by slug alphabetically
			return roomA.Slug < roomB.Slug
		})

		c.HTML(http.StatusOK, "list.html", gin.H{
			"allRooms":        allRooms,
			"isAuthenticated": userId != "",
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
