package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"
)

type Telemetry struct {
	DeviceName string  `json:"device_name"`
	CPU        float64 `json:"cpu"`
	Memory     float64 `json:"memory"`
	Temp       float64 `json:"temperature"`
	Timestamp  string  `json:"timestamp"`
}

func main() {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "device-telemetry",
	})

	http.HandleFunc("/telemetry", func(w http.ResponseWriter, r *http.Request) {
		var t Telemetry
		json.NewDecoder(r.Body).Decode(&t)
		t.Timestamp = time.Now().Format(time.RFC3339)

		data, _ := json.Marshal(t)

		err := writer.WriteMessages(r.Context(),
			kafka.Message{Value: data},
		)
		if err != nil {
			http.Error(w, "Kafka error", 500)
			return
		}

		w.Write([]byte("Telemetry sent"))
	})

	log.Println("Collector API running on :8080")
	http.ListenAndServe(":8080", nil)
}
