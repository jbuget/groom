CREATE TABLE user_starred_rooms (
  user_id VARCHAR(255) NOT NULL,  -- Google user ID
  room_id INT NOT NULL,           -- Reference to rooms.id
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (user_id, room_id),
  FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
);