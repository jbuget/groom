package googleapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"google.golang.org/api/meet/v2"
	"google.golang.org/api/option"
)

// Erreur personnalisée pour l'utilisateur non authentifié
var ErrUserNotAuthenticated = errors.New("user not authenticated")

// GetGoogleAPIClient récupère le client OAuth2 de l'utilisateur authentifié à partir de la session
func GetGoogleAPIClient(c *gin.Context, oauthConfig *oauth2.Config) (*http.Client, error) {
	// Récupérer le token OAuth2 depuis la session
	session := sessions.Default(c)
	tokenJSON := session.Get("token")
	if tokenJSON == nil {
		return nil, ErrUserNotAuthenticated
	}

	// Désérialiser le token JSON en objet oauth2.Token
	var oauthToken oauth2.Token
	err := json.Unmarshal(tokenJSON.([]byte), &oauthToken)
	if err != nil {
		return nil, err
	}

	// Créer et retourner un client OAuth avec le token désérialisé
	client := oauthConfig.Client(context.Background(), &oauthToken)
	return client, nil
}

// GetGoogleMeetRoomInfo récupère les informations d'une Google Meet room, y compris la conférence active
func GetGoogleMeetRoomInfo(c *gin.Context, oauthConfig *oauth2.Config, spaceName string) {
	// Récupérer le client Google API
	client, err := GetGoogleAPIClient(c, oauthConfig)
	if err == ErrUserNotAuthenticated {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Google API client", "details": err.Error()})
		return
	}

	// Créer un service Google Meet
	meetService, err := meet.NewService(c, option.WithHTTPClient(client))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Google Meet service"})
		return
	}

	// Appeler l'API Google Meet pour récupérer les informations de la room
	space, err := meetService.Spaces.Get("spaces/" + spaceName).Do()
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
