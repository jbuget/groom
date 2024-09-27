/**
 * Handlers qui tapent en direct sur l'API Goole
 */
package handlers

import (
	googleapi "groom/internal/google"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/meet/v2"
	"google.golang.org/api/option"
)

func GoogleMeetRoomHandler(c *gin.Context) {
	spaceName := c.Param("name")

	client, exists := c.Get(GoogleClientKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	meetService, err := meet.NewService(c, option.WithHTTPClient(client.(*http.Client)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Google Meet service"})
		return
	}

	space, err := googleapi.GetGoogleMeetRoomInfo(spaceName, meetService)
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
