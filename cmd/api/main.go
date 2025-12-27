// App entry point

package main

import (
	"agnos-gin/config"
	"agnos-gin/internal/domain"
	"agnos-gin/internal/handler"
	"agnos-gin/internal/infrastructure"
	"agnos-gin/internal/middleware"
	"agnos-gin/internal/repository"
	"agnos-gin/internal/service"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// 2. Connect DB
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	// 3. Migrate
	db.AutoMigrate(&domain.Staff{}, &domain.Patient{})

	// 4. Init Layers
	// Repos
	staffRepo := repository.NewStaffRepository(db)
	patientRepo := repository.NewPatientRepository(db)

	// Infrastructure
	hospitalClient := infrastructure.NewHospitalApiClient(cfg.HospitalAPIURL)

	// Services
	authService := service.NewAuthService(staffRepo, cfg.JWTSecret)
	patientService := service.NewPatientService(patientRepo, hospitalClient)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	patientHandler := handler.NewPatientHandler(patientService)

	// 5. Router
	r := gin.Default()

	// Routes
	r.POST("/staff/create", authHandler.Register)
	r.POST("/staff/login", authHandler.Login)

	protected := r.Group("/patient")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		protected.GET("/search", patientHandler.Search)
	}

	// Run
	r.Run(":8080")
}