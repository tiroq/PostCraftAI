package handlers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/tiroq/postcraftai/backend/models"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(func() string {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		s = "default_secret_key_lkmniu*&hvsniMU(A*NVew98fn4gw" // Replace this with a secure secret in production.
	}
	return s
}())

// Claims defines the JWT payload.
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// generateToken creates a JWT token.
func generateToken(username, role string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Signup registers a new user.
func Signup(c *gin.Context) {
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&data); err != nil {
		log.Printf("Signup: bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if _, exists := models.Users[data.Username]; exists {
		log.Printf("Signup: user already exists: %s", data.Username)
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Signup: error hashing password for %s: %v", data.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing password"})
		return
	}
	models.Users[data.Username] = models.User{
		Username:     data.Username,
		PasswordHash: string(hash),
		Role:         "user",
		Allowed:      false, // Must be enabled by admin.
	}
	token, err := generateToken(data.Username, "user")
	if err != nil {
		log.Printf("Signup: token generation failed for %s: %v", data.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}
	log.Printf("Signup: new user created: %s", data.Username)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Login authenticates a user.
func Login(c *gin.Context) {
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&data); err != nil {
		log.Printf("Login: bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	user, exists := models.Users[data.Username]
	if !exists {
		log.Printf("Login: user not found: %s", data.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(data.Password)); err != nil {
		log.Printf("Login: invalid password for user: %s", data.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	token, err := generateToken(user.Username, user.Role)
	if err != nil {
		log.Printf("Login: token generation failed for %s: %v", data.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}
	log.Printf("Login: successful login for user: %s", data.Username)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// AuthMiddleware validates the JWT token.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// AdminMiddleware ensures the caller is an admin.
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if role, exists := c.Get("role"); !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
