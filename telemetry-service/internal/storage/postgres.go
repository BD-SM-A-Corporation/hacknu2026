package storage

import (
	"log"
	"telemetry-service/internal/processor"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

func NewPostgresStorage(dsn string) *PostgresStorage {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("[WARNING] DB connection failed (skipping auto-migration): %v", err)
		return &PostgresStorage{
			db:        nil,
			batchChan: make(chan *TelemetryHistory, 1000),
		}
	}

	db.AutoMigrate(&TelemetryHistory{})
	log.Println("PostgreSQL connection established and schema migrated.")

	ps := &PostgresStorage{
		db:        db,
		batchChan: make(chan *TelemetryHistory, 10000),
	}

	go ps.startBatchProcessor()

	return ps
}

func (ps *PostgresStorage) InsertAsync(state *processor.LocomotiveState) {
	if ps.db == nil {
		return
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

func (ps *PostgresStorage) startBatchProcessor() {
	var batch []*TelemetryHistory
	ticker := time.NewTicker(1 * time.Second)
	batchSize := 100

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

func (ps *PostgresStorage) flush(batch []*TelemetryHistory) {
	if ps.db == nil {
		return
	}
	if err := ps.db.Create(&batch).Error; err != nil {
		log.Printf("Batch insert failed: %v", err)
	}
}
