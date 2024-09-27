package handlers

import (
	"context"
	"encoding/json" // Utilisé pour la désérialisation du token JSON
	"net/http"

	"groom/internal/config"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"                       // Utilisé pour gérer OAuth2
	"golang.org/x/oauth2/google"                // Importer le package meet/v2
	oauth2api "google.golang.org/api/oauth2/v2" // Renommé pour éviter le conflit
	"google.golang.org/api/option"
)

var oauthConfig *oauth2.Config

// Initialisation de la configuration OAuth2
func InitOAuth(cfg config.Config) {
	oauthConfig = &oauth2.Config{
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

// Middleware pour vérifier l'authentification Google
func RequireLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")

		if user == nil {
			// Enregistrer l'URL d'origine dans la session avant de rediriger vers Google OAuth
			session.Set("redirect", c.Request.RequestURI)
			session.Save()

			// Rediriger vers la page de connexion OAuth
			c.Redirect(http.StatusFound, "/auth/login")
			c.Abort()
			return
		}

		c.Next()
	}
}

// Redirige vers Google OAuth
func LoginHandler(c *gin.Context) {
	url := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusFound, url)
}

// Callback après authentification Google
func AuthCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Créer un client avec le token OAuth
	client := oauthConfig.Client(context.Background(), token)

	// Utiliser oauth2api.NewService pour initialiser le service OAuth2
	oauth2Service, err := oauth2api.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create OAuth2 service"})
		return
	}

	// Récupérer les informations utilisateur via l'API Google OAuth2
	userinfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	// Vérifier que l'utilisateur est du domaine inclusion.gouv.fr
	if userinfo.Hd != "inclusion.gouv.fr" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized domain"})
		return
	}

	// Sérialiser le token en JSON pour le stocker dans la session
	tokenJSON, err := json.Marshal(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize token"})
		return
	}

	// Stocker l'utilisateur et le token sérialisé dans la session
	session := sessions.Default(c)
	session.Set("user", userinfo.Email)
	session.Set("token", tokenJSON) // Stocker le token OAuth2 au format JSON
	session.Save()

	// Rediriger l'utilisateur vers l'URL qu'il voulait initialement accéder
	redirect := session.Get("redirect")
	if redirect != nil {
		session.Delete("redirect")
		session.Save()
		c.Redirect(http.StatusFound, redirect.(string))
	} else {
		c.Redirect(http.StatusFound, "/")
	}
}

// Déconnexion
func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/")
}

func BasicAuthMiddleware(username, password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, pass, hasAuth := c.Request.BasicAuth()
		if !hasAuth || user != username || pass != password {
			c.Header("WWW-Authenticate", `Basic realm="restricted"`)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func ApiKeyMiddleware(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestApiKey := c.GetHeader("X-API-KEY")
		if requestApiKey != apiKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
