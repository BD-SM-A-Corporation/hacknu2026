package storage

import (
	"log"
	"telemetry-service/internal/processor"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TelemetryHistory is the DB model for recording historical states.
type TelemetryHistory struct {
	ID           uint   `gorm:"primaryKey"`
	LocomotiveID string `gorm:"index"`
	Speed        float64
	Temperature  float64
	Pressure     float64
	FuelLevel    float64
	Timestamp    time.Time `gorm:"index"`
}

type PostgresStorage struct {
	db        *gorm.DB
	batchChan chan *TelemetryHistory
}

// NewPostgresStorage establishes a PostgreSQL connection and starts the batching worker.
func NewPostgresStorage(dsn string) *PostgresStorage {
	// e.g. dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Almaty"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("[WARNING] DB connection failed (skipping auto-migration): %v", err)
		return &PostgresStorage{
			db:        nil,
			batchChan: make(chan *TelemetryHistory, 1000), // Won't persist but prevents panic
		}
	}

	// Auto migrate the schema (using TimescaleDB in prod is recommended for this setup)
	db.AutoMigrate(&TelemetryHistory{})
	log.Println("PostgreSQL connection established and schema migrated.")

	ps := &PostgresStorage{
		db:        db,
		batchChan: make(chan *TelemetryHistory, 10000), // Buffer for batching
	}

	go ps.startBatchProcessor()

	return ps
}

// InsertAsync queues the record for batch insertion later.
func (ps *PostgresStorage) InsertAsync(state *processor.LocomotiveState) {
	if ps.db == nil {
		return // Ignore if DB is disconnected
	}

	record := &TelemetryHistory{
		LocomotiveID: state.LocomotiveID,
		Speed:        state.Speed,
		Temperature:  state.Temperature,
		Pressure:     state.Pressure,
		FuelLevel:    state.FuelLevel,
		Timestamp:    state.Timestamp,
	}

	ps.batchChan <- record
}

// startBatchProcessor gathers queued records and inserts them in bulk periodically.
func (ps *PostgresStorage) startBatchProcessor() {
	var batch []*TelemetryHistory
	ticker := time.NewTicker(1 * time.Second) // Flush every 1 second
	batchSize := 100                          // Or when 100 records hit

	for {
		select {
		case record := <-ps.batchChan:
			batch = append(batch, record)
			if len(batch) >= batchSize {
				ps.flush(batch)
				batch = nil
			}
		case <-ticker.C:
			if len(batch) > 0 {
				ps.flush(batch)
				batch = nil
			}
		}
	}
}

// flush actually executes the insert query.
func (ps *PostgresStorage) flush(batch []*TelemetryHistory) {
	if ps.db == nil {
		return
	}
	if err := ps.db.Create(&batch).Error; err != nil {
		log.Printf("Batch insert failed: %v", err)
	}
}
