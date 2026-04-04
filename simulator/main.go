package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type LocomotiveData struct {
	LocoID    string `json:"locomotive_id"`
	Series    string `json:"series"` // KZ8A или ТЭ33А
	Timestamp string `json:"timestamp"`

	// Общие параметры
	Speed      float64 `json:"speed"`        // км/ч
	Position   int     `json:"position"`     // Позиция контроллера (0-15 или 0-8)
	PressureTM float64 `json:"oil_pressure"` // Давление в тормозной магистрали (кгс/см²)
	PressureGR float64 `json:"pressure_gr"`  // Давление в главных резервуарах (кгс/см²)

	// Специфика электровоза (KZ8A)
	VoltageCS float64 `json:"voltage_cs,omitempty"` // Напряжение контактной сети (кВ)
	Amperage  float64 `json:"amperage,omitempty"`   // Ток тяги (А)

	// Специфика тепловоза (ТЭ33А)
	FuelLevel  float64 `json:"fuel_perc"`            // Уровень топлива (л)
	EngineTemp float64 `json:"temp"`                 // Температура воды/масла (°C)
	EngineRPM  int     `json:"engine_rpm,omitempty"` // Обороты дизеля

	Status string `json:"status"`
}

func main() {
	url := os.Getenv("WS_URL")
	if url == "" {
		url = "ws://localhost/telemetry/"
	}

	dialer := websocket.DefaultDialer
	var conn *websocket.Conn
	var err error

	// Retry loop for Docker startup sequence
	for i := 0; i < 15; i++ {
		conn, _, err = dialer.Dial(url, nil)
		if err == nil {
			break
		}
		log.Printf("Ожидание сервера %s... (%d/15)", url, i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Ошибка подключения к серверу:", err)
	}
	defer conn.Close()

	// Read pump to handle control messages (pings/pongs) and avoid buffer overflow
	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()

	fmt.Println("Симулятор КТЖ запущен. Отправка данных для KZ8A...")

	for {
		data := generateLocoData("KZ8A-0125")
		payload, _ := json.Marshal(data)

		err := conn.WriteMessage(websocket.TextMessage, payload)
		if err != nil {
			log.Println("Ошибка отправки:", err)
			return
		}

		fmt.Printf("[%s] Скорость: %.1f км/ч | Напряжение: %.1f кВ | Статус: %s\n",
			data.LocoID, data.Speed, data.VoltageCS, data.Status)

		time.Sleep(1 * time.Second) // Имитация частоты опроса датчиков 1 Гц
	}
}

func generateLocoData(id string) LocomotiveData {
	// Имитация случайных изменений
	speed := 60.0 + rand.Float64()*20.0
	voltage := 25.0 + (rand.Float64()*2.0 - 1.0) // Колебания вокруг 25кВ
	status := "Normal"

	// Пример логики критического значения
	if speed > 115 {
		status = "Warning"
	}

	return LocomotiveData{
		LocoID:     id,
		Series:     "KZ8A",
		Timestamp:  time.Now().Format(time.RFC3339),
		Speed:      speed,
		Position:   8,
		PressureTM: 5.2 + rand.Float64()*0.4,
		PressureGR: 9.1,
		VoltageCS:  voltage,
		Amperage:   1200 + rand.Float64()*500,
		FuelLevel:  85.0 - (rand.Float64() * 0.1), // Slow fuel consumption
		EngineTemp: 75.0 + rand.Float64()*10.0,
		Status:     status,
	}
}
