package processor

import (
	"log"
	"math"
	"sync"
)

// WorkerPool manages incoming raw telemetry bytes, processes them concurrently,
// and outputs clean LocomotiveStates to a specific channel.
type WorkerPool struct {
	incoming    chan []byte
	outgoing    chan *LocomotiveState
	transformer DataTransformer
	numWorkers  int
	wg          sync.WaitGroup
	prevStates  map[string]*LocomotiveState
	mu          sync.RWMutex
}

// NewWorkerPool initializes the pool.
func NewWorkerPool(workers int, transformer DataTransformer) *WorkerPool {
	return &WorkerPool{
		incoming:    make(chan []byte, 10000), // Buffer for high-load bursts
		outgoing:    make(chan *LocomotiveState, 10000),
		transformer: transformer,
		numWorkers:  workers,
		prevStates:  make(map[string]*LocomotiveState),
	}
}

// Start launches the goroutines for data processing.
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// computeDistance calculates the approximate distance between two points in kilometers.
func computeDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371 // Earth radius in km
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

// worker converts raw bytes to parsed structure.
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	log.Printf("Worker %d started", id)

	for rawPayload := range wp.incoming {
		state, err := wp.transformer.Transform(rawPayload)
		if err != nil {
			log.Printf("[Worker %d] Failed to parse: %v", id, err)
			continue
		}

		// Basic boundary validation
		state.GpsCorrupted = false
		if state.Lat < 40.0 || state.Lat > 56.0 || state.Lng < 46.0 || state.Lng > 88.0 {
			state.GpsCorrupted = true
		}

		wp.mu.Lock()
		if prevState, exists := wp.prevStates[state.LocomotiveID]; exists {
			timeDiff := state.Timestamp.Sub(prevState.Timestamp).Seconds()
			if timeDiff > 0 {
				// 1. Detect GPS Jumps (Bug)
				// If distance moved > 5km in a short interval (e.g. < 60s), it's likely a GPS bug.
				dist := computeDistance(prevState.Lat, prevState.Lng, state.Lat, state.Lng)
				if dist > 5.0 && timeDiff < 60.0 {
					log.Printf("[Worker %d] GPS Jump Detected for %s: %.2f km in %.2fs", id, state.LocomotiveID, dist, timeDiff)
					state.GpsCorrupted = true
					state.IsAnomaly = true
				}

				// 2. Speed change > 40 km/h per second
				if (math.Abs(state.Speed-prevState.Speed) / timeDiff) > 40.0 {
					state.IsAnomaly = true
				}

				// 3. Fuel jump/drop > 5% means likely anomaly
				if math.Abs(state.FuelLevel-prevState.FuelLevel) > 5.0 {
					state.IsAnomaly = true
				}
			}
		}
		// Save for next comparison (only if GPS is NOT corrupted, otherwise we'd compare against a bad point next)
		if !state.GpsCorrupted {
			wp.prevStates[state.LocomotiveID] = state
		}
		wp.mu.Unlock()

		wp.outgoing <- state
	}
}

// Submit enqueues data for processing.
func (wp *WorkerPool) Submit(data []byte) {
	wp.incoming <- data
}

// Output returns the channel containing processed states.
func (wp *WorkerPool) Output() <-chan *LocomotiveState {
	return wp.outgoing
}
