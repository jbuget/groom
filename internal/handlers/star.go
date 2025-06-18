package handlers

import (
	"database/sql"
	"groom/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// StarRoomHandler handles the starring of a room by a user
func StarRoomHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get room ID from URL
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		// Get user ID from session
		session := sessions.Default(c)
		userID := session.Get("user")
		if userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Star the room
		err = models.StarRoom(db, userID.(string), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to star room"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Room starred successfully"})
	}
}

// UnstarRoomHandler handles the unstarring of a room by a user
func UnstarRoomHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get room ID from URL
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		// Get user ID from session
		session := sessions.Default(c)
		userID := session.Get("user")
		if userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Unstar the room
		err = models.UnstarRoom(db, userID.(string), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unstar room"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Room unstarred successfully"})
	}
}

// ToggleStarRoomHandler handles toggling the star status of a room
func ToggleStarRoomHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get room ID from URL
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		// Get user ID from session
		session := sessions.Default(c)
		userID := session.Get("user")
		if userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Check if room is already starred
		isStarred, err := models.IsRoomStarredByUser(db, userID.(string), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check star status"})
			return
		}

		// Toggle star status
		var message string
		if isStarred {
			err = models.UnstarRoom(db, userID.(string), id)
			message = "Room unstarred successfully"
		} else {
			err = models.StarRoom(db, userID.(string), id)
			message = "Room starred successfully"
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle star status"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": message,
			"starred": !isStarred,
		})
	}
}