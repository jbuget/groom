package models

import (
	"database/sql"
	"time"
)

type Room struct {
	ID        int       `json:"id"`
	Slug      string    `json:"slug"`
	SpaceID   string    `json:"space_id"`
	IsStarred bool      `json:"is_starred"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetRoomByID(db *sql.DB, id int) (*Room, error) {
	row := db.QueryRow("SELECT id, slug, space_id, created_at, updated_at FROM rooms WHERE id = $1", id)

	var room Room
	room.IsStarred = false // Default value

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
	room.IsStarred = false // Default value

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
		room.IsStarred = false // Default value
		if err := rows.Scan(&room.ID, &room.Slug, &room.SpaceID, &room.CreatedAt, &room.UpdatedAt); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

// GetAllRoomsWithStarStatus gets all rooms with starred status for a user
// Rooms are ordered by starred status (starred first), then by slug
func GetAllRoomsWithStarStatus(db *sql.DB, userID string) ([]Room, error) {
	var rooms []Room
	
	query := `
		SELECT r.id, r.slug, r.space_id, r.created_at, r.updated_at,
		       EXISTS(SELECT 1 FROM user_starred_rooms sr WHERE sr.room_id = r.id AND sr.user_id = $1) AS is_starred
		FROM rooms r
		ORDER BY 
		    EXISTS(SELECT 1 FROM user_starred_rooms sr WHERE sr.room_id = r.id AND sr.user_id = $1) DESC,
		    r.slug ASC
	`
	
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var room Room
		if err := rows.Scan(&room.ID, &room.Slug, &room.SpaceID, &room.CreatedAt, &room.UpdatedAt, &room.IsStarred); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, err
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

