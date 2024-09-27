package main

import (
	"groom/internal/config"
	"groom/internal/db"
	"groom/internal/handlers"
	"log"

	"github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
)

func main() {
	// Charger la configuration
	cfg := config.LoadConfig()

	// Connexion à la base de données
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer database.Close()

	// Appliquer les migrations au démarrage
	if err := db.RunMigrations(database, "./migrations"); err != nil {
		log.Fatalf("Could not run migrations: %v\n", err)
	}

    // Initialisation d'OAuth
	handlers.InitOAuth(cfg)

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
    r.GET("/", handlers.RequireLogin(), handlers.ListRoomsHTMLHandler(database))

	// Route pour rediriger avec un slug
	r.GET("/:slug", handlers.RedirectHandler(database))

	// Routes API protégées par clé d'API
	api := r.Group("/api", handlers.ApiKeyMiddleware(cfg.APIKey))
	{
        api.GET("/rooms", handlers.ListRoomsJSONHandler(database))
        api.POST("/rooms", handlers.CreateRoomHandler(database))
        api.PUT("/rooms/:id", handlers.UpdateRoomHandler(database))
        api.DELETE("/rooms/:id", handlers.DeleteRoomHandler(database))
	}

	// Démarrer le serveur
	log.Printf("Server started at %s:%s", cfg.Host, cfg.Port)
	if err := r.Run(cfg.Host + ":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
