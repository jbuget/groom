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

	// Serves static files
	r.Static("/static", "./static")

	// Routes pour l'authentification Google
	r.GET("/auth/login", handlers.LoginHandler)
	r.GET("/auth/callback", handlers.AuthCallbackHandler(cfg.GoogleWorkspaceDomain))
	r.GET("/auth/logout", handlers.LogoutHandler)

	// Protected routes (by "X-API-TOKEN" HTTP header)
	api := r.Group("/api", handlers.ApiKeyMiddleware(cfg.APIKey))
	{
		api.GET("/rooms", handlers.ListRoomsJSONHandler(db.Database))
		api.GET("/room-status", handlers.GetRoomStatusHandler(db.Database, googleapi.MeetService))
		api.POST("/rooms", handlers.CreateRoomHandler(db.Database, googleapi.MeetService))
		api.PUT("/rooms/:id", handlers.UpdateRoomHandler(db.Database))
		api.DELETE("/rooms/:id", handlers.DeleteRoomHandler(db.Database))
	}
	
	// User routes (require login)
	user := r.Group("/api/user", handlers.RequireLogin())
	{
		user.POST("/rooms/:id/star", handlers.StarRoomHandler(db.Database))
		user.DELETE("/rooms/:id/star", handlers.UnstarRoomHandler(db.Database))
		user.POST("/rooms/:id/toggle-star", handlers.ToggleStarRoomHandler(db.Database))
	}

	// System routes
	r.GET("/healthz", handlers.HealthzHandler(db.Database, googleapi.MeetService))

	// Open routes
	r.GET("/", handlers.RequireLogin(), handlers.ListRoomsHTMLHandler(db.Database, googleapi.MeetService))
	r.GET("/status", handlers.RequireLogin(), handlers.GetRoomStatusHandler(db.Database, googleapi.MeetService))
	r.GET("/:slug", handlers.RedirectHandler(db.Database, googleapi.MeetService))

	// Démarrer le serveur
	log.Printf("Server started at %s:%s", cfg.Host, cfg.Port)
	if err := r.Run(cfg.Host + ":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
