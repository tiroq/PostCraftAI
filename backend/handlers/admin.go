package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tiroq/postcraftai/backend/models"
)

// AdminEnableUser allows an admin to enable a user with access expiration (in minutes).
func AdminEnableUser(c *gin.Context) {
	var data struct {
		Username  string `json:"username"`
		ExpiresIn int    `json:"expires_in"` // in minutes
	}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	user, exists := models.Users[data.Username]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if data.ExpiresIn <= 0 {
		data.ExpiresIn = 10080 // default to 7 days in minutes.
	}
	user.Allowed = true
	user.AccessExpiresAt = time.Now().Add(time.Duration(data.ExpiresIn) * time.Minute)
	models.Users[data.Username] = user
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("User %s enabled until %s", data.Username, user.AccessExpiresAt.Format(time.RFC1123)),
	})
}

// AdminListUsers returns the list of registered users.
func AdminListUsers(c *gin.Context) {
	var list []gin.H
	for _, u := range models.Users {
		list = append(list, gin.H{
			"username":         u.Username,
			"role":             u.Role,
			"allowed":          u.Allowed,
			"access_expiresAt": u.AccessExpiresAt.Format(time.RFC1123),
		})
	}
	c.JSON(http.StatusOK, list)
}

// AdminUpdateRateLimit allows admin to update the global OpenAI rate limit.
func AdminUpdateRateLimit(c *gin.Context) {
	var data struct {
		RateLimit int `json:"rate_limit"`
	}
	if err := c.BindJSON(&data); err != nil || data.RateLimit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rate_limit value"})
		return
	}
	models.OpenAIRateLimit = data.RateLimit
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Rate limit updated to %d req/min", data.RateLimit)})
}

// AdminUpdateExpiration allows an admin to update an already enabled user's expiration.
func AdminUpdateExpiration(c *gin.Context) {
	var data struct {
		Username  string `json:"username"`
		ExpiresIn int    `json:"expires_in"` // in minutes
	}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	user, exists := models.Users[data.Username]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if data.ExpiresIn <= 0 {
		data.ExpiresIn = 10080 // default to 7 days in minutes.
	}
	user.AccessExpiresAt = time.Now().Add(time.Duration(data.ExpiresIn) * time.Minute)
	models.Users[data.Username] = user
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("User %s expiration updated until %s", data.Username, user.AccessExpiresAt.Format(time.RFC1123)),
	})
}

// AdminRequestStats returns the logged request activities.
func AdminRequestStats(c *gin.Context) {
	c.JSON(http.StatusOK, models.RequestLogs())
}
