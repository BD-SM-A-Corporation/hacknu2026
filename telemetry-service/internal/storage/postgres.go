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
	LocomotiveID string `gorm:"index:idx_loco_time"`
	Speed        float64
	Temperature  float64
	Pressure     float64
	FuelLevel    float64
	Lat          float64
	Lng          float64
	IsAnomaly    bool
	Timestamp    time.Time `gorm:"index:idx_loco_time;index:idx_timestamp"`
}

func (TelemetryHistory) TableName() string {
	return "telemetry_logs"
}

type PostgresStorage struct {
	db *gorm.DB
}

func NewPostgresStorage(dsn string) *PostgresStorage {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("[WARNING] DB connection failed (skipping auto-migration): %v", err)
		return &PostgresStorage{
			db: nil,
		}
	}

	db.AutoMigrate(&TelemetryHistory{})
	log.Println("PostgreSQL connection established and schema migrated.")

	return &PostgresStorage{
		db: db,
	}
}

func (ps *PostgresStorage) InsertSync(state *processor.LocomotiveState) error {
	if ps.db == nil {
		return nil
	}

	record := &TelemetryHistory{
		LocomotiveID: state.LocomotiveID,
		Speed:        state.Speed,
		Temperature:  state.Temperature,
		Pressure:     state.Pressure,
		FuelLevel:    state.FuelLevel,
		Lat:          state.Lat,
		Lng:          state.Lng,
		IsAnomaly:    state.IsAnomaly,
		Timestamp:    state.Timestamp,
	}

	return ps.db.Create(record).Error
}
