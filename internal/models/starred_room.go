package models

import (
	"database/sql"
	"time"
)

// StarredRoom represents a room that has been starred by a user
type StarredRoom struct {
	UserID    string
	RoomID    int64
	CreatedAt time.Time
}

// StarRoom marks a room as starred by a user
func StarRoom(db *sql.DB, userID string, roomID int64) error {
	query := `
		INSERT INTO user_starred_rooms (user_id, room_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, room_id) DO NOTHING
	`
	_, err := db.Exec(query, userID, roomID)
	return err
}

// UnstarRoom removes a room from a user's starred rooms
func UnstarRoom(db *sql.DB, userID string, roomID int64) error {
	query := `
		DELETE FROM user_starred_rooms
		WHERE user_id = $1 AND room_id = $2
	`
	_, err := db.Exec(query, userID, roomID)
	return err
}

// IsRoomStarredByUser checks if a room is starred by a user
func IsRoomStarredByUser(db *sql.DB, userID string, roomID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM user_starred_rooms
			WHERE user_id = $1 AND room_id = $2
		)
	`
	var isStarred bool
	err := db.QueryRow(query, userID, roomID).Scan(&isStarred)
	return isStarred, err
}

// GetStarredRoomIDsByUser returns all room IDs starred by a user
func GetStarredRoomIDsByUser(db *sql.DB, userID string) ([]int64, error) {
	query := `
		SELECT room_id FROM user_starred_rooms
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roomIDs []int64
	for rows.Next() {
		var roomID int64
		if err := rows.Scan(&roomID); err != nil {
			return nil, err
		}
		roomIDs = append(roomIDs, roomID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roomIDs, nil
}