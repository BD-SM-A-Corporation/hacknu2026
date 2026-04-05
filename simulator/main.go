package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type LocomotiveData struct {
	LocoID     string  `json:"locomotive_id"`
	Series     string  `json:"series"`
	Timestamp  string  `json:"timestamp"`
	Speed      float64 `json:"speed"`
	Position   int     `json:"position"`
	PressureTM float64 `json:"oil_pressure"`
	PressureGR float64 `json:"pressure_gr"`
	VoltageCS  float64 `json:"voltage_cs,omitempty"`
	Amperage   float64 `json:"amperage,omitempty"`
	FuelLevel  float64 `json:"fuel_perc"`
	EngineTemp float64 `json:"temp"`
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
	Status     string  `json:"status"`
}

func main() {
	url := os.Getenv("WS_URL")
	if url == "" {
		url = "ws://telemetry-go:8080/"
	}

	locomotives := []struct {
		id     string
		series string
	}{
		{"KZ8A-0125", "KZ8A"},
		{"KZ8A-0130", "KZ8A"},
		{"TE33A-5501", "TE33A"},
	}

	var wg sync.WaitGroup
	for _, l := range locomotives {
		wg.Add(1)
		go simulateLoco(url, l.id, l.series, &wg)
	}
	wg.Wait()
}

func simulateLoco(url, id, series string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Initial State
	currentFuel := 85.0
	if series == "KZ8A" {
		currentFuel = 0
	}

	lat, lng := 51.1693, 71.4490 // Astana
	if id == "KZ8A-0130" {
		lat, lng = 43.2220, 76.8512 // Almaty
	} else if id == "KZ8A-0125" {
		lat, lng = 49.8018, 73.1021 // Karaganda
	}

	fmt.Printf(">>> Запущен %s (%s)\n", id, series)

	for {
		dialer := websocket.DefaultDialer
		conn, _, err := dialer.Dial(url, nil)
		if err != nil {
			log.Printf("[%s] Ошибка подключения: %v. Повтор через 5с...", id, err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Printf("[%s] Соединение установлено", id)

		// Message reader to handle close events
		done := make(chan struct{})
		go func() {
			defer close(done)
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					log.Printf("[%s] Ошибка чтения (соединение разорвано): %v", id, err)
					return
				}
			}
		}()

		// Simulation Loop
	loop:
		for {
			select {
			case <-done:
				break loop
			default:
				chance := rand.Float64()
				isSilent := false

				data := generateBaseData(id, series, &currentFuel, &lat, &lng)

				if chance < 0.05 {
					anomalyType := rand.Intn(5)
					switch anomalyType {
					case 0:
						fmt.Printf("⚠️ [%s] ОБРЫВ СВЯЗИ\n", id)
						isSilent = true
					case 1:
						fmt.Printf("⚠️ [%s] СКАЧОК СКОРОСТИ\n", id)
						data.Speed = 185.5
						data.Status = "Critical"
					case 2:
						fmt.Printf("⚠️ [%s] ПЕРЕГРЕВ\n", id)
						data.EngineTemp = 115.0
						data.Status = "Warning"
					case 3:
						if series == "TE33A" {
							fmt.Printf("⚠️ [%s] СЛИВ ТОПЛИВА (-15%%)\n", id)
							currentFuel -= 15.0
							if currentFuel < 0 {
								currentFuel = 0
							}
							data.FuelLevel = currentFuel
							data.Status = "Critical"
						}
					case 4:
						fmt.Printf("⚠️ [%s] GPS JUMP\n", id)
						lat += 0.5
						data.Lat = lat
						data.Status = "Critical"
					}
				}

				if !isSilent {
					payload, _ := json.Marshal(data)
					if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
						log.Printf("[%s] Ошибка записи: %v", id, err)
						conn.Close()
						break loop
					}
				} else {
					time.Sleep(15 * time.Second)
				}

				time.Sleep(1 * time.Second)
			}
		}
		conn.Close()
	}
}

func generateBaseData(id, series string, fuel *float64, lat *float64, lng *float64) LocomotiveData {
	speed := 40.0 + rand.Float64()*40.0

	// Расход топлива
	if series == "TE33A" {
		*fuel -= 0.005
		if *fuel < 0 {
			*fuel = 0
		}
	}

	// Движение (очень упрощенно: едем на Северо-Восток)
	dist := speed / 3600.0 / 111.0
	*lat += dist * 0.7
	*lng += dist * 0.7

	data := LocomotiveData{
		LocoID:     id,
		Series:     series,
		Timestamp:  time.Now().Format(time.RFC3339),
		Speed:      speed,
		Position:   rand.Intn(9),
		PressureTM: 5.2 + (rand.Float64() * 0.1),
		PressureGR: 9.0 + (rand.Float64() * 0.2),
		EngineTemp: 70.0 + rand.Float64()*10.0,
		FuelLevel:  *fuel,
		Lat:        *lat,
		Lng:        *lng,
		Status:     "Normal",
	}

	if series == "KZ8A" {
		data.VoltageCS = 25.0 + (rand.Float64() * 0.5)
		data.Amperage = 1000 + rand.Float64()*200
	}

	return data
}
