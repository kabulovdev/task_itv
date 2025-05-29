package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/itv_task/models"
	"gorm.io/gorm"
)

type MovieHandler struct {
	DB *gorm.DB
}

func NewMovieHandler(db *gorm.DB) *MovieHandler {
	return &MovieHandler{DB: db}
}

// CreateMovie godoc
// @Summary Create a new movie
// @Description Create a new movie record
// @Tags movies
// @Accept json
// @Produce json
// @Param movie body models.Movie true "Movie to create"
// @Success 201 {object} models.Movie
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Security BearerAuth
// @Router /api/movies [post]
func (h *MovieHandler) CreateMovie(c *gin.Context) {
	var movie models.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	if err := h.DB.Create(&movie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, movie)
}

// GetMovies godoc
// @Summary Get all movies
// @Description Retrieve all movies with optional search and pagination
// @Tags movies
// @Produce json
// @Param title query string false "Search by title (partial match)"
// @Param year query int false "Search by year"
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Page size (default 10)"
// @Success 200 {array} models.Movie
// @Failure 500 {object} handlers.ErrorResponse
// @Security BearerAuth
// @Router /api/movies [get]
func (h *MovieHandler) GetMovies(c *gin.Context) {
	var movies []models.Movie

	title := c.Query("title")
	yearStr := c.Query("year")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	db := h.DB.Model(&models.Movie{})
	if title != "" {
		db = db.Where("title ILIKE ?", "%"+title+"%")
	}
	if yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err == nil {
			db = db.Where("year = ?", year)
		}
	}
	if err := db.Offset(offset).Limit(limit).Find(&movies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, movies)
}

// GetMovieByID godoc
// @Summary Get a movie by ID
// @Description Retrieve a specific movie by ID
// @Tags movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} models.Movie
// @Failure 404 {object} handlers.ErrorResponse
// @Security BearerAuth
// @Router /api/movies/{id} [get]
func (h *MovieHandler) GetMovieByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var movie models.Movie
	if err := h.DB.First(&movie, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Movie not found"})
		return
	}
	c.JSON(http.StatusOK, movie)
}

// UpdateMovie godoc
// @Summary Update a movie
// @Description Update an existing movie record
// @Tags movies
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Param movie body models.Movie true "Movie to update"
// @Success 200 {object} models.Movie
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Security BearerAuth
// @Router /api/movies/{id} [put]
func (h *MovieHandler) UpdateMovie(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var movie models.Movie
	if err := h.DB.First(&movie, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Movie not found"})
		return
	}
	var input models.Movie
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	movie.Title = input.Title
	movie.Director = input.Director
	movie.Year = input.Year
	movie.Plot = input.Plot
	if err := h.DB.Save(&movie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, movie)
}

// DeleteMovie godoc
// @Summary Delete a movie
// @Description Delete a movie by ID
// @Tags movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} SuccessResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Security BearerAuth
// @Router /api/movies/{id} [delete]
func (h *MovieHandler) DeleteMovie(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.DB.Delete(&models.Movie{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{Message: "Movie deleted"})
}
