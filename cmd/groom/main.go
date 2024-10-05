package main

import (
	"groom/internal/config"
	"groom/internal/db"
	googleapi "groom/internal/google"
	"groom/internal/handlers"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// Charger la configuration
	cfg := config.LoadConfig()

	// Initialisation de la base de donées et exécution automatique des migrations
	db.InitDatabase(cfg.DatabaseURL)
	databaseName := "postgres"
	if err := db.RunMigrations(cfg.DatabaseMigrationPath, databaseName); err != nil {
		log.Fatalf("Could not run migrations: %v\n", err)
	}
	defer db.Database.Close()

	// Initialisation des composants Google (OAuth utilisateur ou compte de services, clients d'APIs, etc.)
	googleapi.InitUserOAuth(cfg)
	googleapi.InitServiceAccountServices(cfg)

	// Création du routeur Gin
	r := gin.Default()

	// Configuration de la session
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// Chargement des templates HTML
	r.LoadHTMLGlob("templates/*")

	// Routes pour l'authentification Google
	r.GET("/auth/login", handlers.LoginHandler)
	r.GET("/auth/callback", handlers.AuthCallbackHandler)
	r.GET("/auth/logout", handlers.LogoutHandler)

	// Protection des routes par OAuth
	r.GET("/", handlers.RequireLogin(), handlers.ListRoomsHTMLHandler(db.Database))

	// Route pour rediriger avec un slug
	r.GET("/:slug", handlers.RedirectHandler(db.Database))

	// Route pour avoir des infos sur une room Google Meet
	r.GET("/spaces/:name", handlers.RequireLogin(), handlers.GoogleMeetRoomHandler(googleapi.MeetService))

	// Routes API protégées par clé d'API
	api := r.Group("/api", handlers.ApiKeyMiddleware(cfg.APIKey))
	{
		api.GET("/rooms", handlers.ListRoomsJSONHandler(db.Database))
		api.POST("/rooms", handlers.CreateRoomHandler(db.Database, googleapi.MeetService))
		api.PUT("/rooms/:id", handlers.UpdateRoomHandler(db.Database))
		api.DELETE("/rooms/:id", handlers.DeleteRoomHandler(db.Database))
	}

	// Démarrer le serveur
	log.Printf("Server started at %s:%s", cfg.Host, cfg.Port)
	if err := r.Run(cfg.Host + ":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
