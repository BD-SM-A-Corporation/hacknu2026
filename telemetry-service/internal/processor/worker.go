package processor

import (
	"log"
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
}

// NewWorkerPool initializes the pool.
func NewWorkerPool(workers int, transformer DataTransformer) *WorkerPool {
	return &WorkerPool{
		incoming:    make(chan []byte, 10000), // Buffer for high-load bursts
		outgoing:    make(chan *LocomotiveState, 10000),
		transformer: transformer,
		numWorkers:  workers,
	}
}

// Start launches the goroutines for data processing.
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
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

		// Further rules, validation could be placed here before sending it out.

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
