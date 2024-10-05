package googleapi

import (
	"context"
	"encoding/json"
	"groom/internal/config"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google" // Importer le package meet/v2
	"golang.org/x/oauth2/jwt"
	meet "google.golang.org/api/meet/v2"
	"google.golang.org/api/option"
)

var UserOAuthConfig *oauth2.Config
var ServiceAccountOAuthConfig *jwt.Config
var MeetService *meet.Service

func InitUserOAuth(cfg config.Config) {
	UserOAuthConfig = &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func InitServiceAccountServices(cfg config.Config) {
	ctx := context.Background()

	serviceAccountFile := "./service_account.json"

	credentialsJSON, err := os.ReadFile(serviceAccountFile)
	if err != nil {
		log.Fatalf("Unable to read service account file: %v", err)
	}

	// Désérialiser le fichier JSON
	var credentials struct {
		ClientEmail string `json:"client_email"`
		PrivateKey  string `json:"private_key"`
		TokenURL    string `json:"token_uri"`
	}
	if err := json.Unmarshal(credentialsJSON, &credentials); err != nil {
		log.Fatalf("Unable to unmarshal service account JSON: %v", err)
	}

	// Spécifiez l'utilisateur à impersonner (un utilisateur du domaine Google Workspace)
	impersonatedUser := "collectif@inclusion.gouv.fr" // L'utilisateur que vous voulez impersonner

	// Configurer le compte de service pour agir en tant qu'utilisateur avec délégation
	ServiceAccountOAuthConfig := &jwt.Config{
		Email:      credentials.ClientEmail,
		PrivateKey: []byte(credentials.PrivateKey),
		Scopes: []string{
			meet.MeetingsSpaceCreatedScope,
			meet.MeetingsSpaceReadonlyScope,
		},
		TokenURL: credentials.TokenURL,
		Subject:  impersonatedUser, // Spécifiez l'utilisateur pour l'impersonation
	}
	client := ServiceAccountOAuthConfig.Client(context.Background())

	MeetService, err = meet.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

}
