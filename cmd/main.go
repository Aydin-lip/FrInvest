package main

import (
	"log"
	"recruitment-api/config"
	"recruitment-api/internal/controller"
	"recruitment-api/internal/email"
	"recruitment-api/internal/middleware"
	"recruitment-api/internal/models"
	"recruitment-api/internal/repository"
	"recruitment-api/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	config.Load()

	// Connect to DB
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Auto migrate
	if err := db.AutoMigrate(&models.User{}, &models.VerificationCode{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	// Repositories
	userRepo := repository.NewUserRepository(db)
	verificationRepo := repository.NewVerificationRepository(db)

	// Services
	jwtService := service.NewJWTService()
	mailer := email.NewMailer()
	authService := service.NewAuthService(userRepo, verificationRepo, jwtService, mailer)
	userService := service.NewUserService(userRepo)

	// Controllers
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)

	// Router
	r := gin.Default()

	// Auth routes (public)
	auth := r.Group("/api/auth")
	{
		auth.POST("/send-code", authController.SendCode)
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	// User routes (protected)
	users := r.Group("/api/users")
	users.Use(middleware.AuthMiddleware(jwtService, userRepo))
	{
		users.GET("", userController.GetAll)
		users.PUT("/status", userController.UpdateStatus)
		users.GET("/statistics", userController.GetStatistics)
	}

	port := ":" + config.AppConfig.AppPort
	log.Printf("Server starting on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
