package websocket

import (
	"encoding/json"
	"telemetry-service/internal/processor"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the processor to relay outward.
	broadcast chan *processor.LocomotiveState

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Where to send data incoming from clients? (Usually sensors would use separate endpoints, but maybe they share WS)
	IncomingTelemetry chan []byte
}

func NewHub() *Hub {
	return &Hub{
		broadcast:         make(chan *processor.LocomotiveState, 5000),
		register:          make(chan *Client),
		unregister:        make(chan *Client),
		clients:           make(map[*Client]bool),
		IncomingTelemetry: make(chan []byte, 10000),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			// Package the state properly into a JSON wrapper if frontend requires it.
			// Currently frontend expects: { type: 'telemetry', data: { locomotiveId, speed, ... } }
			payload := map[string]interface{}{
				"type": "telemetry",
				"data": message,
			}

			rawBytes, err := json.Marshal(payload)
			if err != nil {
				continue
			}

			for client := range h.clients {
				select {
				case client.send <- rawBytes:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// BroadcastState puts the state into the broadcast queue.
func (h *Hub) BroadcastState(state *processor.LocomotiveState) {
	h.broadcast <- state
}
