package handlers

import (
	googleapi "groom/internal/google"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/meet/v2"
)

func GoogleMeetRoomHandler(meetService *meet.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		spaceName := c.Param("name")

		space, err := googleapi.MeetService.Spaces.Get("spaces/" + spaceName).Do()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve room info", "details": err.Error()})
			return
		}

		// Vérifier s'il y a une conférence active
		activeConference := false
		if space.ActiveConference != nil {
			activeConference = true
		}

		// Retourner les informations de la room avec la conférence active (boolean)
		c.JSON(http.StatusOK, gin.H{
			"space":            space,
			"activeConference": activeConference,
		})
	}
}
