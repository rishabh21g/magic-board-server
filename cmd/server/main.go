package main

import (
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rishabh21g/magic-board/internal/config"
	"github.com/rishabh21g/magic-board/internal/game"
	"github.com/rishabh21g/magic-board/internal/middleware"
	"github.com/rishabh21g/magic-board/internal/store"
	"github.com/rishabh21g/magic-board/internal/ws"
)

func main() {
	cfg := config.LoadConfig()
	// initialize Redis client and store
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddress,
	})
	redisStore := store.NewRedisStore(redisClient)

	// initialize rate limiter
	rl := middleware.NewRateLimiter(redisClient, 1, time.Second)

	// game service and WebSocket handler initialization
	gameService := game.NewService(redisStore)

	// initialize WebSocket hub and handler
	hub := ws.NewHub()
	go hub.Run()

	handler := ws.NewHandler(hub, gameService, rl, redisStore)

	// handle WebSocket connections at /ws endpoint
	http.HandleFunc("/ws", handler.ServeWs)

	//health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// start HTTP server
	log.Printf("Server started on port %s", cfg.Port)
	err := http.ListenAndServe(":"+cfg.Port, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
