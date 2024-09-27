/**
 * Handlers qui tapent en direct sur l'API Goole
 */
package handlers

import (
	googleapi "groom/internal/google"

	"github.com/gin-gonic/gin"
)

func GoogleMeetRoomHandler(c *gin.Context) {
	spaceName := c.Param("name")
	googleapi.GetGoogleMeetRoomInfo(c, oauthConfig, spaceName)
}
