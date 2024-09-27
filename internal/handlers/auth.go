package handlers

import (
	"context"
	"net/http"

	"groom/internal/config"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2" // Utilisé pour gérer OAuth2
	"golang.org/x/oauth2/google"
	oauth2api "google.golang.org/api/oauth2/v2" // Renommé pour éviter le conflit
)

var oauthConfig *oauth2.Config

// Initialisation de la configuration OAuth2
func InitOAuth(cfg config.Config) {
	oauthConfig = &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

// Middleware pour vérifier l'authentification Google
func RequireLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")

		if user == nil {
			// Redirige vers Google OAuth si l'utilisateur n'est pas connecté
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

	client := oauthConfig.Client(context.Background(), token)
	service, err := oauth2api.New(client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create OAuth2 service"})
		return
	}

	userinfo, err := service.Userinfo.Get().Do()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	// Vérifier que l'utilisateur est du domaine inclusion.gouv.fr
	if userinfo.Hd != "inclusion.gouv.fr" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized domain"})
		return
	}

	// Stocker l'utilisateur dans la session
	session := sessions.Default(c)
	session.Set("user", userinfo.Email)
	session.Save()

	c.Redirect(http.StatusFound, "/")
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
