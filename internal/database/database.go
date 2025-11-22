package database

import (
	"context"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Service interface {
	Health() map[string]string
	Close() error
	GetDB() *sqlx.DB
}

type service struct {
	db *sqlx.DB
}

func New(dsn string) Service {
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &service{
		db: db,
	}
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"
	stats["open_connections"] = fmt.Sprintf("%d", s.db.Stats().OpenConnections)

	return stats
}

func (s *service) Close() error {
	return s.db.Close()
}

func (s *service) GetDB() *sqlx.DB {
	return s.db
}
