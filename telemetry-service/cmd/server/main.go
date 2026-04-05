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
	log.Println("Starting Locomotive Telemetry Service...")

	hub := websocket.NewHub()
	go hub.Run()

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := broker.NewRedisBroker(redisAddr, "", 0)

	// change for production
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=telemetry port=5432 sslmode=disable"
	}
	dbStorage := storage.NewPostgresStorage(dsn)

	transformer := &processor.DefaultTransformer{}
	pool := processor.NewWorkerPool(10, transformer)
	pool.Start()

	go func() {
		for {
			select {
			case rawPayload := <-hub.IncomingTelemetry:
				pool.Submit(rawPayload)

			case state := <-pool.Output():
				if err := rdb.UpdateHeartbeat(context.Background(), state.LocomotiveID); err != nil {
					log.Printf("Heartbeat error: %v", err)
				}

				if !state.GpsCorrupted {
					hub.BroadcastState(state)
				}

				if err := dbStorage.InsertSync(state); err != nil {
					log.Printf("DB insert error: %v", err)
				} else {
					if err := rdb.Publish(context.Background(), "locomotive-updates", state); err != nil {
						log.Printf("Redis publish error: %v", err)
					}
				}
			}
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})

	log.Println("Telemetry service listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Fatal server error: %v", err)
	}
}
