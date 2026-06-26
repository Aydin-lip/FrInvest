package main

import (
	"log"
	"recruitment-api/config"
	"recruitment-api/internal/controller"
	"recruitment-api/internal/email"
	"recruitment-api/internal/models"
	"recruitment-api/internal/repository"
	"recruitment-api/internal/service"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.VerificationToken{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	verificationRepo := repository.NewVerificationRepository(db)

	mailer := email.NewMailer()
	authService := service.NewAuthService(userRepo, verificationRepo, mailer)
	userService := service.NewUserService(userRepo)

	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)
	dashboardController := controller.NewDashboardController(userService)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/send-verification", authController.SendVerification)
		auth.GET("/verify-email", authController.VerifyEmail)
	}

	users := r.Group("/api/users")
	{
		users.GET("", userController.GetAll)
		users.PATCH("/status", userController.UpdateStatus)
	}

	dashboard := r.Group("/api/dashboard")
	{
		dashboard.GET("/status-percentages", dashboardController.GetStatusPercentages)
	}

	port := ":" + config.AppConfig.AppPort
	log.Printf("Server starting on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
