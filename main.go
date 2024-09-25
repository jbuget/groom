package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"

    _ "github.com/jackc/pgx/v4/stdlib"
)

var db *sql.DB

// Room struct to hold room data
type Room struct {
    ID     int    `json:"id"`
    Slug   string `json:"slug"`
    MeetID string `json:"meet_id"`
}

// Middleware pour protéger l'accès avec un mot de passe
func basicAuth(next http.HandlerFunc, username, password string) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user, pass, ok := r.BasicAuth()
        if !ok || user != username || pass != password {
            w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next(w, r)
    })
}

// Middleware pour vérifier la clé d'API
func apiKeyMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        apiKey := os.Getenv("USHR_API_KEY")
        if apiKey == "" {
            http.Error(w, "API key not configured", http.StatusInternalServerError)
            return
        }

        // Récupérer la clé d'API envoyée dans les en-têtes
        requestApiKey := r.Header.Get("X-API-KEY")
        if requestApiKey != apiKey {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Si la clé est valide, continuer vers le prochain handler
        next(w, r)
    })
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

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return rooms, nil
}

// Handler pour afficher la liste des rooms
func listRoomsHTMLHandler(w http.ResponseWriter, r *http.Request) {
    rooms, err := getAllRooms()
    if err != nil {
        http.Error(w, "Unable to retrieve rooms", http.StatusInternalServerError)
        return
    }

    // Génération de la réponse HTML simple
    fmt.Fprintf(w, "<h1>Liste des salles</h1>")
    fmt.Fprintf(w, "<ul>")
    for _, room := range rooms {
        roomURL := fmt.Sprintf("https://meet.google.com/%s", room.MeetID)
        fmt.Fprintf(w, `<li><strong>%s</strong>: <a href="%s" target="_blank">%s</a></li>`, room.Slug, roomURL, roomURL)
    }
    fmt.Fprintf(w, "</ul>")
}

// Handler pour lister les rooms en format JSON
func listRoomsJSONHandler(w http.ResponseWriter, r *http.Request) {
    rooms, err := getAllRooms()
    if err != nil {
        http.Error(w, "Unable to retrieve rooms", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(rooms)
}

// Handler pour rediriger en fonction du slug
func redirectHandler(w http.ResponseWriter, r *http.Request) {
    slug := strings.TrimPrefix(r.URL.Path, "/") // Récupère le slug sans le "/"
    meetID, err := getMeetIDFromSlug(slug)
    if err != nil {
        http.Error(w, "Room not found", http.StatusNotFound)
        return
    }

    targetURL := fmt.Sprintf("https://meet.google.com/%s", meetID)
    http.Redirect(w, r, targetURL, http.StatusFound)
}

// Fonction pour récupérer l'ID de la room Google Meet à partir du slug
func getMeetIDFromSlug(slug string) (string, error) {
    var meetID string
    query := "SELECT meet_id FROM rooms WHERE slug = $1"
    err := db.QueryRow(query, slug).Scan(&meetID)
    if err != nil {
        return "", err
    }
    return meetID, nil
}

// Handler POST pour ajouter une room
func createRoomHandler(w http.ResponseWriter, r *http.Request) {
    var room Room
    err := json.NewDecoder(r.Body).Decode(&room)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    query := "INSERT INTO rooms (slug, meet_id) VALUES ($1, $2) RETURNING id"
    err = db.QueryRow(query, room.Slug, room.MeetID).Scan(&room.ID)
    if err != nil {
        http.Error(w, "Error inserting room", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(room)
}

// Handler PUT pour modifier une room
func updateRoomHandler(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/api/rooms/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid room ID", http.StatusBadRequest)
        return
    }

    var room Room
    err = json.NewDecoder(r.Body).Decode(&room)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    query := "UPDATE rooms SET slug = $1, meet_id = $2 WHERE id = $3"
    _, err = db.Exec(query, room.Slug, room.MeetID, id)
    if err != nil {
        http.Error(w, "Error updating room", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Room updated successfully")
}

// Handler DELETE pour supprimer une room
func deleteRoomHandler(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/api/rooms/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid room ID", http.StatusBadRequest)
        return
    }

    query := "DELETE FROM rooms WHERE id = $1"
    _, err = db.Exec(query, id)
    if err != nil {
        http.Error(w, "Error deleting room", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Room deleted successfully")
}

func main() {
    var err error
    dbURL := os.Getenv("DATABASE_URL") // Ex: postgres://user:pass@localhost:5432/dbname
    db, err = sql.Open("pgx", dbURL)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v\n", err)
    }
    defer db.Close()

    // Récupérer les variables d'environnement HOST et PORT
    host := os.Getenv("HOST")
    if host == "" {
        host = "0.0.0.0" // Par défaut, écouter sur toutes les interfaces
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "3000" // Par défaut, écouter sur le port 3080
    }
	
    // Définir les identifiants pour l'accès protégé
    username := os.Getenv("BASIC_AUTH_LOGIN")		// Remplacez par votre nom d'utilisateur
    password := os.Getenv("BASIC_AUTH_PASSWORD")	// Remplacez par votre mot de passe

    // Route par défaut pour lister les rooms ou rediriger, protégée par mot de passe
    http.HandleFunc("/", basicAuth(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/" {
            // Si l'URL est la racine, lister les rooms
            listRoomsHTMLHandler(w, r)
        } else {
            // Sinon, traiter comme un slug
            redirectHandler(w, r)
        }
    }, username, password))

    // Routes sous /api pour les opérations JSON protégées par clé d'API
    http.HandleFunc("/api/rooms", apiKeyMiddleware(func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            listRoomsJSONHandler(w, r)
        case http.MethodPost:
            createRoomHandler(w, r)
        default:
            http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        }
    }))

    http.HandleFunc("/api/rooms/", apiKeyMiddleware(func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodPut:
            updateRoomHandler(w, r)
        case http.MethodDelete:
            deleteRoomHandler(w, r)
        default:
            http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        }
    }))

    // Lancer le serveur sur l'hôte et le port définis
    addr := fmt.Sprintf("%s:%s", host, port)
    fmt.Printf("Server started at %s\n", addr)
    log.Fatal(http.ListenAndServe(addr, nil))
}