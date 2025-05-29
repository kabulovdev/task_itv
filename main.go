// @title ITV movies Gateway v1
// @version 1.0
// @description This is a ITV movies Gateway server.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"

	"github.com/itv_task/db"
	_ "github.com/itv_task/docs"
	"github.com/itv_task/handlers"
	"github.com/itv_task/middleware"
	"github.com/itv_task/models"
	"gorm.io/gorm"
)

func main() {
	fx.New(
		fx.Provide(db.Connect),
		fx.Provide(NewRouter),
		fx.Provide(handlers.NewMovieHandler),
		fx.Provide(handlers.NewUserHandler),
		fx.Invoke(MigrateDB),
		fx.Invoke(RegisterRoutes),
		fx.Invoke(StartServer),
	).Run()
}

func NewRouter() *gin.Engine {
	r := gin.Default()
	return r
}

func MigrateDB(db *gorm.DB) {
	models.Migrate(db)
	db.AutoMigrate(&models.User{})
}

func RegisterRoutes(r *gin.Engine, h *handlers.MovieHandler, uh *handlers.UserHandler) {
	api := r.Group("/api")
	{
		api.POST("/register", uh.Register)
		api.POST("/login", uh.Login)
		movies := api.Group("/movies")
		movies.Use(middleware.JWTAuthMiddleware())
		{
			movies.POST("", middleware.AdminOnly(), h.CreateMovie)
			movies.GET("", middleware.UserOnly(), h.GetMovies)
			movies.GET(":id", middleware.UserOnly(), h.GetMovieByID)
			movies.PUT(":id", middleware.AdminOnly(), h.UpdateMovie)
			movies.DELETE(":id", middleware.AdminOnly(), h.DeleteMovie)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func StartServer(r *gin.Engine) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
