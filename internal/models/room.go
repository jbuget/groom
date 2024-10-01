package models

import (
	"database/sql"
)

type Room struct {
	ID      int    `json:"id"`
	Slug    string `json:"slug"`
	SpaceID string `json:"space_id"`
}

// Fonction pour récupérer toutes les rooms
func GetAllRooms(db *sql.DB) ([]Room, error) {
	var rooms []Room
	rows, err := db.Query("SELECT id, slug, space_id FROM rooms ORDER BY slug ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var room Room
		if err := rows.Scan(&room.ID, &room.Slug, &room.SpaceID); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

// Fonction pour créer une nouvelle room
func CreateRoom(db *sql.DB, room Room) (int, error) {
	var id int
	query := "INSERT INTO rooms (slug, space_id) VALUES ($1, $2) RETURNING id"
	err := db.QueryRow(query, room.Slug, room.SpaceID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Fonction pour mettre à jour une room
func UpdateRoom(db *sql.DB, id int, room Room) error {
	query := "UPDATE rooms SET slug = $1, space_id = $2 WHERE id = $3"
	_, err := db.Exec(query, room.Slug, room.SpaceID, id)
	return err
}

// Fonction pour supprimer une room
func DeleteRoom(db *sql.DB, id int) error {
	query := "DELETE FROM rooms WHERE id = $1"
	_, err := db.Exec(query, id)
	return err
}

// Récupérer l'ID de la room Google Meet à partir du slug
func GetSpaceIDFromSlug(db *sql.DB, slug string) (string, error) {
	var spaceID string
	query := "SELECT space_id FROM rooms WHERE slug = $1"
	err := db.QueryRow(query, slug).Scan(&spaceID)
	if err != nil {
		return "", err
	}
	return spaceID, nil
}
