package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/itv_task/models"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}


// Register godoc
// @Summary Register a new user
// @Description Register a new user with username, email, password, and role
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User registration info"
// @Success 201 {object} models.User
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 409 {object} handlers.ErrorResponse
// @Router /api/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	if user.Role == "" {
		user.Role = "user"
	}
	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, ErrorResponse{Error: "Username or email already exists"})
		return
	}
	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary Login
// @Description Authenticate user and get JWT token (login with username and password only)
// @Tags auth
// @Accept json
// @Produce json
// @Param username query string true "Username"
// @Param password query string true "Password"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "username and password are required as query parameters"})
		return
	}
	var user models.User
	if err := h.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid username or password"})
		return
	}
	if !user.CheckPassword(password) {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid username or password"})
		return
	}
	// Generate JWT with user ID and role
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": "Bearer " + tokenString})
}
