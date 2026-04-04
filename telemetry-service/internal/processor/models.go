package processor

import (
	"encoding/json"
	"time"
)

// RawTelemetry represents incoming unstable data (from different sources).
// We parse it loosely and transform it to LocomotiveState.
type RawTelemetry struct {
	LocomotiveID string  `json:"locomotive_id"`
	Speed        float64 `json:"speed"`
	Temp         float64 `json:"temp"`
	OilPressure  float64 `json:"oil_pressure"`
	FuelPerc     float64 `json:"fuel_perc"`
}

// LocomotiveState is the validated and unified structure for the frontend and DB.
type LocomotiveState struct {
	LocomotiveID string    `json:"locomotiveId"`
	Speed        float64   `json:"speed"`
	Temperature  float64   `json:"temperature"`
	Pressure     float64   `json:"pressure"`
	FuelLevel    float64   `json:"fuelLevel"`
	IsAnomaly    bool      `json:"is_anomaly"`
	Timestamp    time.Time `json:"timestamp"`
}

// DataTransformer is an interface to allow different algorithms to shape incoming data.
type DataTransformer interface {
	Transform(raw []byte) (*LocomotiveState, error)
}

// DefaultTransformer parses standard JSON structure.
type DefaultTransformer struct{}

func (t *DefaultTransformer) Transform(raw []byte) (*LocomotiveState, error) {
	var rt RawTelemetry
	if err := json.Unmarshal(raw, &rt); err != nil {
		return nil, err
	}

	// Add potential data cleaning rules here if necessary.
	return &LocomotiveState{
		LocomotiveID: rt.LocomotiveID,
		Speed:        rt.Speed,
		Temperature:  rt.Temp,
		Pressure:     rt.OilPressure,
		FuelLevel:    rt.FuelPerc,
		Timestamp:    time.Now(),
	}, nil
}
