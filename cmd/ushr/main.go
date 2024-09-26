package main

import (
	"log"
	"ushr/internal/config"
	"ushr/internal/db"
	"ushr/internal/handlers"

	"github.com/gin-gonic/gin"
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

	// Création du routeur Gin
	r := gin.Default()

	// Chargement des templates HTML
	r.LoadHTMLGlob("templates/*")

	// Routes HTML avec authentification basique
	r.GET("/", handlers.BasicAuthMiddleware(cfg.Username, cfg.Password), handlers.ListRoomsHTMLHandler(database))

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
