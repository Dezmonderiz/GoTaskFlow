package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"GoTaskFlow/internal/cache"
	"GoTaskFlow/internal/config"
	"GoTaskFlow/internal/handler"
	"GoTaskFlow/internal/repository"
	"GoTaskFlow/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	db, err := openDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	redisClient, err := openRedis(cfg)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer redisClient.Close()

	taskRepository := repository.NewPostgresTaskRepository(db)
	statsCache := cache.NewRedisStatsCache(redisClient, time.Duration(cfg.StatsCacheTTL)*time.Second)
	taskService := service.NewTaskService(taskRepository, statsCache)
	router := handler.NewRouter(taskService)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func openDatabase(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			return nil, errors.Join(err, closeErr)
		}
		return nil, err
	}

	return db, nil
}

func openRedis(cfg config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		if closeErr := client.Close(); closeErr != nil {
			return nil, errors.Join(err, closeErr)
		}
		return nil, err
	}

	return client, nil
}
