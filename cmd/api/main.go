package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/echo/internal/config"
	"github.com/echo/internal/database"
	"github.com/echo/internal/handlers"
	"github.com/echo/internal/models"
	"github.com/echo/internal/pkg/jwt"
	"github.com/echo/internal/repository"
	router "github.com/echo/internal/router"
	"github.com/echo/internal/service"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	cfg, _ := config.Load()

	s := cfg.Server
	serverAddr := fmt.Sprintf(":%s", s.Port)

	db, err := database.NewDatabaseConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to Database: %v", err)
	}
	log.Println("✓ Database connection established")

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("✓ Database migration completed")

	redis, err := database.NewRedisConnection(&cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redis.Close()
	log.Println("✓ Redis connection established")

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	userRepo := repository.NewUserRepository(db)

	jwtManager := jwt.NewManager(cfg.Jwt.Secret, cfg.Jwt.AccessExpiry, cfg.Jwt.RefreshExpiry)

	//Service
	authService := service.NewAuthService(userRepo, *jwtManager)
	userService := service.NewUserService(userRepo)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	router.Setup(e, authHandler, userHandler, jwtManager)

	go func() {
		if err := e.Start(serverAddr); err != nil {
			e.Logger.Error("failed to start server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	log.Println("Server exited")
}
