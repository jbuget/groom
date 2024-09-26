package main

import (
    "database/sql"
    "github.com/gin-gonic/gin"
    "log"
    "net/http"
    "os"
    "strconv"

    _ "github.com/jackc/pgx/v4/stdlib"
)

var db *sql.DB

// Room struct to hold room data
type Room struct {
    ID     int    `json:"id"`
    Slug   string `json:"slug"`
    MeetID string `json:"meet_id"`
}

// Middleware pour l'authentification basique
func basicAuthMiddleware(username, password string) gin.HandlerFunc {
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

// Middleware pour vérifier la clé d'API
func apiKeyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        apiKey := os.Getenv("USHR_API_KEY")
        if apiKey == "" {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "API key not configured"})
            c.Abort()
            return
        }

        requestApiKey := c.GetHeader("X-API-KEY")
        if requestApiKey != apiKey {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        c.Next()
    }
}

// Fonction pour récupérer toutes les rooms depuis la base de données
func getAllRooms() ([]Room, error) {
    var rooms []Room
    rows, err := db.Query("SELECT id, slug, meet_id FROM rooms")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var room Room
        if err := rows.Scan(&room.ID, &room.Slug, &room.MeetID); err != nil {
            return nil, err
        }
        rooms = append(rooms, room)
    }

    return rooms, nil
}

// Handler pour afficher la liste des rooms en HTML
func listRoomsHTMLHandler(c *gin.Context) {
    rooms, err := getAllRooms()
    if err != nil {
        c.String(http.StatusInternalServerError, "Unable to retrieve rooms")
        return
    }

    c.HTML(http.StatusOK, "list.html", gin.H{
        "rooms": rooms,
    })
}

// Handler pour lister les rooms en JSON
func listRoomsJSONHandler(c *gin.Context) {
    rooms, err := getAllRooms()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve rooms"})
        return
    }

    c.JSON(http.StatusOK, rooms)
}

// Handler pour rediriger en fonction du slug
func redirectHandler(c *gin.Context) {
    slug := c.Param("slug")
    var meetID string
    query := "SELECT meet_id FROM rooms WHERE slug = $1"
    err := db.QueryRow(query, slug).Scan(&meetID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
        return
    }

    c.Redirect(http.StatusFound, "https://meet.google.com/"+meetID)
}

// Handler POST pour ajouter une room
func createRoomHandler(c *gin.Context) {
    var room Room
    if err := c.BindJSON(&room); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    query := "INSERT INTO rooms (slug, meet_id) VALUES ($1, $2) RETURNING id"
    err := db.QueryRow(query, room.Slug, room.MeetID).Scan(&room.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting room"})
        return
    }

    c.JSON(http.StatusCreated, room)
}

// Handler PUT pour modifier une room
func updateRoomHandler(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
        return
    }

    var room Room
    if err := c.BindJSON(&room); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    query := "UPDATE rooms SET slug = $1, meet_id = $2 WHERE id = $3"
    _, err = db.Exec(query, room.Slug, room.MeetID, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating room"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Room updated successfully"})
}

// Handler DELETE pour supprimer une room
func deleteRoomHandler(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
        return
    }

    query := "DELETE FROM rooms WHERE id = $1"
    _, err = db.Exec(query, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting room"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Room deleted successfully"})
}

func main() {
    // Connexion à la base de données
    var err error
    dbURL := os.Getenv("DATABASE_URL")
    db, err = sql.Open("pgx", dbURL)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v\n", err)
    }
    defer db.Close()

    // Récupérer les variables d'environnement HOST et PORT
    host := os.Getenv("HOST")
    if host == "" {
        host = "0.0.0.0"
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }

    // Définir les identifiants pour l'authentification basique
    username := os.Getenv("BASIC_AUTH_LOGIN")
    password := os.Getenv("BASIC_AUTH_PASSWORD")

    // Création du routeur Gin
    r := gin.Default()
    
    // Chargement des templates HTML
    r.LoadHTMLGlob("templates/*")

    // Appliquer basicAuth uniquement à la route racine "/"
    r.GET("/", basicAuthMiddleware(username, password), listRoomsHTMLHandler)

    // La route "/:slug" redirige sans authentification
    r.GET("/:slug", redirectHandler)

    // Routes sous /api pour les opérations JSON protégées par clé d'API
    api := r.Group("/api", apiKeyMiddleware())
    {
        api.GET("/rooms", listRoomsJSONHandler)
        api.POST("/rooms", createRoomHandler)
        api.PUT("/rooms/:id", updateRoomHandler)
        api.DELETE("/rooms/:id", deleteRoomHandler)
    }

    // Démarrer le serveur
    addr := host + ":" + port
    log.Printf("Server started at %s", addr)
    if err := r.Run(addr); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}