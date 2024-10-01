package handlers

import (
	"database/sql"
	"groom/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler pour rediriger en fonction du slug
func RedirectHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")
		spaceID, err := models.GetSpaceIDFromSlug(db, slug)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}

		// Rediriger vers la room Google Meet correspondante
		c.Redirect(http.StatusFound, "https://meet.google.com/"+spaceID)
	}
}

// Handler pour afficher la liste des rooms en HTML
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
