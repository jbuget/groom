package googleapi

import (
	"context"
	"groom/internal/config"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google" // Importer le package meet/v2
	meet "google.golang.org/api/meet/v2"
)

var OAuthConfig *oauth2.Config

// Initialisation de la configuration OAuth2
func InitOAuth(cfg config.Config) {
	OAuthConfig = &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/meetings.space.created",  // Scope pour Google Meet API
			"https://www.googleapis.com/auth/meetings.space.readonly", // Scope pour Google Meet API
		},
		Endpoint: google.Endpoint,
	}
}

func GetGoogleAPIClientFromOAuth2Token(oauthToken *oauth2.Token) (*http.Client, error) {
	client := OAuthConfig.Client(context.Background(), oauthToken)
	return client, nil
}

func GetGoogleMeetRoomInfo(spaceName string, meetService *meet.Service) (*meet.Space, error) {
	return meetService.Spaces.Get("spaces/" + spaceName).Do()
}
