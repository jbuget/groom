package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func BasicAuthMiddleware(username, password string) gin.HandlerFunc {
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

func ApiKeyMiddleware(apiKey string) gin.HandlerFunc {
    return func(c *gin.Context) {
        requestApiKey := c.GetHeader("X-API-KEY")
        if requestApiKey != apiKey {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }
        c.Next()
    }
}