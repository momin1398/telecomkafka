package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type Telemetry struct {
	DeviceName string  `json:"device_name"`
	CPU        float64 `json:"cpu"`
	Memory     float64 `json:"memory"`
	Temp       float64 `json:"temperature"`
}

func main() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "device-telemetry",
		GroupID: "alarm-service",
	})

	for {
		msg, _ := reader.ReadMessage(context.Background())

		var t Telemetry
		json.Unmarshal(msg.Value, &t)

		if t.CPU > 80 {
			log.Printf("🚨 CPU ALARM on %s: %.2f%%", t.DeviceName, t.CPU)
		}
		if t.Temp > 70 {
			log.Printf("🔥 TEMP ALARM on %s: %.2f°C", t.DeviceName, t.Temp)
		}
	}
}
