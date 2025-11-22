package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"template/internal/auth"
	"template/internal/config"
	"template/internal/database"
	"template/internal/email"
	"template/internal/redis"
	"template/internal/server"
	"template/internal/telemetry"
	"template/internal/user"
	"template/internal/validator"
)

func main() {
	// 1. Load Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 2. Init Telemetry
	shutdownOTel, err := telemetry.InitTracer(context.Background())
	if err != nil {
		log.Printf("failed to init telemetry: %v", err)
	}
	defer func() {
		if err := shutdownOTel(context.Background()); err != nil {
			log.Printf("failed to shutdown telemetry: %v", err)
		}
	}()

	// 3. Init DB
	db := database.New(cfg.DB.DSN)
	defer db.Close()

	// 4. Init Redis
	redisClient := redis.New(cfg.Redis.Addr)

	// 5. Init Validator
	v := validator.New()

	// 6. Init Repos & Services
	emailSender := email.NewSender(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Sender)
	userRepo := user.NewRepository(db.GetDB())
	userService := user.NewService(userRepo, cfg.JWTSecret, emailSender, cfg.FrontendHost)

	// 7. Init Handlers
	authHandler := auth.NewHandler(userService, v)
	userHandler := user.NewHandler(userRepo)

	// 8. Init Server
	srv := server.NewServer(cfg, db, redisClient, authHandler, userHandler)

	// 9. Start Server (Graceful Shutdown)
	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server forced to shutdown: %v", err)
	}

	log.Println("server exited properly")
}
