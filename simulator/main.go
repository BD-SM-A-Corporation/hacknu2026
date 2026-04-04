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
		url = "ws://telemetry-go:8080/" // Имя сервиса из docker-compose
	}

	// Список локомотивов для симуляции
	locomotives := []struct {
		id     string
		series string
	}{
		{"KZ8A-0125", "KZ8A"},
		{"KZ8A-0130", "KZ8A"},
		{"TE33A-5501", "TЭ33A"},
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

	dialer := websocket.DefaultDialer
	var conn *websocket.Conn
	var err error

	// Retry loop
	for i := 0; i < 20; i++ {
		conn, _, err = dialer.Dial(url, nil)
		if err == nil {
			break
		}
		log.Printf("[%s] Подключение к %s... (%d/20)", id, url, i+1)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Printf("[%s] Не удалось подключиться: %v", id, err)
		return
	}
	defer conn.Close()

	// Read pump to handle control messages (pings/pongs) and avoid server write timeout
	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()

	// Состояние локомотива для плавных изменений
	currentFuel := 85.0
	if series == "KZ8A" {
		currentFuel = 0 // У электровозов нет бака в литрах в данном контексте
	}

	// Astana coordinates
	lat := 51.169392
	lng := 71.449074
	if id == "KZ8A-0130" {
		lat = 43.222015 // Almaty
		lng = 76.851248
	} else if id == "KZ8A-0125" {
		lat = 49.801867 // Karaganda
		lng = 73.102143
	}

	fmt.Printf(">>> Запущен поток для %s (%s)\n", id, series)

	for {
		// --- ГЕНЕРАЦИЯ EDGE CASES (АНОМАЛИЙ) ---
		chance := rand.Float64()
		isSilent := false

		data := generateBaseData(id, series, &currentFuel, &lat, &lng)

		if chance < 0.05 { // 5% шанс на аномалию
			anomalyType := rand.Intn(5)
			switch anomalyType {
			case 0:
				fmt.Printf("⚠️ [%s] ИМИТАЦИЯ ОБРЫВА СВЯЗИ (30с)\n", id)
				isSilent = true
			case 1:
				fmt.Printf("⚠️ [%s] ИМИТАЦИЯ СКАЧКА СКОРОСТИ\n", id)
				data.Speed = 185.5
				data.Status = "Critical"
			case 2:
				fmt.Printf("⚠️ [%s] ИМИТАЦИЯ ПЕРЕГРЕВА\n", id)
				data.EngineTemp = 115.0
				data.Status = "Warning"
			case 3:
				if series == "TЭ33A" {
					fmt.Printf("⚠️ [%s] ИМИТАЦИЯ СЛИВА ТОПЛИВА\n", id)
					currentFuel -= 15.0
					if currentFuel < 0 {
						currentFuel = 0
					}
					data.FuelLevel = currentFuel
					data.Status = "Critical"
				}
			case 4:
				fmt.Printf("⚠️ [%s] ИМИТАЦИЯ ПРЫЖКА КООРДИНАТ (GPS СБОЙ)\n", id)
				lat += 1.0 // jump ~111km
				data.Lat = lat
				data.Status = "Critical"
			}
		}

		if !isSilent {
			payload, _ := json.Marshal(data)
			if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				log.Printf("[%s] Ошибка: %v", id, err)
				return
			}
		} else {
			// Если имитируем "Выкл", просто ждем и ничего не шлем
			time.Sleep(30 * time.Second)
		}

		time.Sleep(1 * time.Second)
	}
}

func generateBaseData(id, series string, fuel *float64, lat *float64, lng *float64) LocomotiveData {
	speed := 40.0 + rand.Float64()*40.0
	temp := 70.0 + rand.Float64()*15.0

	// Медленный расход топлива
	if series == "TЭ33A" {
		*fuel -= 0.001
		if *fuel < 0 {
			*fuel = 0
		}
	}

	// Физическое движение координат
	importMath := true
	if importMath {
		distanceKm := speed / 3600.0 // за 1 секунду
		*lat += (distanceKm / 111.0)
		*lng += (distanceKm / 111.0) // simplified to constant factor for brevity
	}

	data := LocomotiveData{
		LocoID:     id,
		Series:     series,
		Timestamp:  time.Now().Format(time.RFC3339),
		Speed:      speed,
		Position:   rand.Intn(9),
		PressureTM: 5.2 + (rand.Float64() * 0.2),
		PressureGR: 9.0 + (rand.Float64() * 0.5),
		EngineTemp: temp,
		FuelLevel:  *fuel,
		Lat:        *lat,
		Lng:        *lng,
		Status:     "Normal",
	}

	if series == "KZ8A" {
		data.VoltageCS = 25.0 + (rand.Float64()*2.0 - 1.0)
		data.Amperage = 800 + rand.Float64()*600
	}

	return data
}
