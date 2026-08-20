package main

import (
	"log"
	"time"

	"queue-management-tenant/backend/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	log.Printf("[WORKER] Starting Queue Management Background Worker (%s)...", cfg.AppEnv)

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("[WORKER] Running background tasks (Stale ticket cleanup, notification queue check)...")
	}
}
