package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"telemetry-service/internal/broker"
	"telemetry-service/internal/processor"
	"telemetry-service/internal/storage"
	"telemetry-service/internal/websocket"
)

func main() {
	// 1. Initialize core dependencies
	log.Println("Starting Locomotive Telemetry Service...")

	// Create hub
	hub := websocket.NewHub()
	go hub.Run()

	// Create Redis Broker (connecting to local redis server on 6379 by default)
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := broker.NewRedisBroker(redisAddr, "", 0)

	// Create PostgreSQL Storage Component
	// DSN format (change for production)
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=telemetry port=5432 sslmode=disable"
	}
	dbStorage := storage.NewPostgresStorage(dsn)

	// Create Data Processor Pool
	transformer := &processor.DefaultTransformer{}
	// Using 10 workers for concurrent processing
	pool := processor.NewWorkerPool(10, transformer)
	pool.Start()

	// 2. Main data pipeline orchestration
	go func() {
		for {
			select {
			// Listen to incoming telemetry from websocket/device connections
			case rawPayload := <-hub.IncomingTelemetry:
				// Forward to worker pool
				pool.Submit(rawPayload)

			// Listen to processed telemetry states output
			case state := <-pool.Output():
				// 0. Update Heartbeat in Redis
				if err := rdb.UpdateHeartbeat(context.Background(), state.LocomotiveID); err != nil {
					log.Printf("Heartbeat error: %v", err)
				}

				// 1. Broadcast to all active Frontend WebSocket subscribers
				hub.BroadcastState(state)

				// 2. Synchronous write row to PostgreSQL
				if err := dbStorage.InsertSync(state); err != nil {
					log.Printf("DB insert error: %v", err)
				} else {
					// 3. Publish to Redis for Laravel backend to react AFTER DB write
					if err := rdb.Publish(context.Background(), "locomotive-updates", state); err != nil {
						log.Printf("Redis publish error: %v", err)
					}
				}
			}
		}
	}()

	// 3. HTTP Endpoints Setup
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})

	// Add CORS mapping if Laravel app sits on different port
	// but standard usage within Homestead or similar handles CORS at API Gateway level.

	log.Println("Telemetry service listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Fatal server error: %v", err)
	}
}
