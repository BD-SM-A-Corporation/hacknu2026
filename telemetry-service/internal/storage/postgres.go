package storage

import (
	"log"
	"telemetry-service/internal/processor"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	db     *gorm.DB
	buffer chan *TelemetryHistory
}

func NewPostgresStorage(dsn string) *PostgresStorage {
	// GORM Tuning: SkipDefaultTransaction and PrepareStmt for performance
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Warn), // Only log slow/error SQL
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		log.Printf("[WARNING] DB connection failed (skipping auto-migration): %v", err)
		return &PostgresStorage{
			db: nil,
		}
	}

	// Connection Pool Tuning
	sqldb, err := db.DB()
	if err == nil {
		sqldb.SetMaxOpenConns(25)
		sqldb.SetMaxIdleConns(10)
		sqldb.SetConnMaxLifetime(5 * time.Minute)
	}

	// We still run AutoMigrate but Laravel migration will handle indexes more explicitly
	db.AutoMigrate(&TelemetryHistory{})
	log.Println("PostgreSQL connection established and schema migrated.")

	ps := &PostgresStorage{
		db:     db,
		buffer: make(chan *TelemetryHistory, 1000), // Buffer for up to 1000 records
	}

	go ps.startBatchWorker()

	return ps
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

	// Non-blocking send to buffer
	select {
	case ps.buffer <- record:
		return nil
	default:
		log.Printf("[WARNING] Telemetry buffer full, dropping record for %s", state.LocomotiveID)
		return nil
	}
}

func (ps *PostgresStorage) startBatchWorker() {
	ticker := time.NewTicker(500 * time.Millisecond) // Flush every 500ms
	defer ticker.Stop()

	var batch []*TelemetryHistory
	batchSize := 100

	for {
		select {
		case record := <-ps.buffer:
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
	if ps.db == nil || len(batch) == 0 {
		return
	}

	// Use GORM CreateInBatches for high efficiency
	err := ps.db.CreateInBatches(batch, len(batch)).Error
	if err != nil {
		log.Printf("[ERROR] Failed to flush telemetry batch: %v", err)
	}
}
