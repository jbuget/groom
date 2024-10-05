package models

import (
	"database/sql"
	"time"
)

type Room struct {
	ID        int       `json:"id"`
	Slug      string    `json:"slug"`
	SpaceID   string    `json:"space_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetRoomByID(db *sql.DB, id int) (*Room, error) {
	row := db.QueryRow("SELECT id, slug, space_id, created_at, updated_at FROM rooms WHERE id = $1", id)

	var room Room

	err := row.Scan(&room.ID, &room.Slug, &room.SpaceID, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &room, nil
}

func GetRoomBySlug(db *sql.DB, slug string) (*Room, error) {
	row := db.QueryRow("SELECT id, slug, space_id, created_at, updated_at FROM rooms WHERE slug = $1", slug)

	var room Room

	err := row.Scan(&room.ID, &room.Slug, &room.SpaceID, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &room, nil
}

func GetAllRooms(db *sql.DB) ([]Room, error) {
	var rooms []Room
	rows, err := db.Query("SELECT id, slug, space_id, created_at, updated_at FROM rooms ORDER BY slug ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var room Room
		if err := rows.Scan(&room.ID, &room.Slug, &room.SpaceID, &room.CreatedAt, &room.UpdatedAt); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func CreateRoom(db *sql.DB, room Room) (*Room, error) {
	query := `
		INSERT INTO rooms (slug, space_id, created_at, updated_at) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, slug, space_id, created_at, updated_at`

	err := db.QueryRow(query, room.Slug, room.SpaceID, time.Now(), time.Now()).
		Scan(&room.ID, &room.Slug, &room.SpaceID, &room.CreatedAt, &room.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &room, nil
}

func UpdateRoom(db *sql.DB, room Room) error {
	query := `
		UPDATE rooms 
		SET slug = $1, space_id = $2, updated_at = $3 
		WHERE id = $4`
	_, err := db.Exec(query, room.Slug, room.SpaceID, time.Now(), room.ID)
	return err
}

func DeleteRoom(db *sql.DB, id int) error {
	query := "DELETE FROM rooms WHERE id = $1"
	_, err := db.Exec(query, id)
	return err
}

func GetSpaceIDFromSlug(db *sql.DB, slug string) (string, error) {
	var spaceID string
	query := "SELECT space_id FROM rooms WHERE slug = $1"
	err := db.QueryRow(query, slug).Scan(&spaceID)
	if err != nil {
		return "", err
	}
	return spaceID, nil
}
