package handlers

import (
	"database/sql"
	"groom/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/meet/v2"
)

// GET /
func ListRoomsHTMLHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rooms, err := models.GetAllRooms(db)
		if err != nil {
			c.String(http.StatusInternalServerError, "Unable to retrieve rooms")
			return
		}

		// Render du template list.html avec la liste des rooms
		c.HTML(http.StatusOK, "list.html", gin.H{
			"rooms": rooms,
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

		if room.SpaceID == "" {
			space, err := meetService.Spaces.Create(&meet.Space{}).Do()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating Google Meet space"})
				return
			}
			room.SpaceID = space.MeetingCode
			room.UpdatedAt = time.Now()
			models.UpdateRoom(db, *room)
		}

		// TODO check GMeet space validity

		// Rediriger vers la room Google Meet correspondante
		c.Redirect(http.StatusFound, "https://meet.google.com/"+room.SpaceID)
	}
}
