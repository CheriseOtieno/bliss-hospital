package handlers

import (
	"context"
	"net/http"
	"os"
	"time"

	"bliss-backend/db"
	"bliss-backend/middleware"
	"bliss-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role := req.Role
	if role == "" {
		role = "patient"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	var user models.User
	query := `
		INSERT INTO users (full_name, email, phone, password_hash, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING user_id, full_name, email, phone, role, is_active, created_at
	`
	err = db.DB.QueryRow(context.Background(), query,
		req.FullName, req.Email, req.Phone, string(hash), role,
	).Scan(&user.UserID, &user.FullName, &user.Email, &user.Phone, &user.Role, &user.IsActive, &user.CreatedAt)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	token, err := generateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{Token: token, User: user})
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	var passwordHash string
	query := `
		SELECT user_id, full_name, email, phone, role, is_active, created_at, password_hash
		FROM users WHERE email = $1 AND is_active = TRUE
	`
	err := db.DB.QueryRow(context.Background(), query, req.Email).Scan(
		&user.UserID, &user.FullName, &user.Email, &user.Phone,
		&user.Role, &user.IsActive, &user.CreatedAt, &passwordHash,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := generateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{Token: token, User: user})
}

func Me(c *gin.Context) {
	userID := c.GetString("user_id")

	var user models.User
	query := `SELECT user_id, full_name, email, phone, role, is_active, created_at FROM users WHERE user_id = $1`
	err := db.DB.QueryRow(context.Background(), query, userID).Scan(
		&user.UserID, &user.FullName, &user.Email, &user.Phone, &user.Role, &user.IsActive, &user.CreatedAt,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func generateToken(user models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	claims := &middleware.Claims{
		UserID: user.UserID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

