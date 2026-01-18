package model

type TelemetryEvent struct {
	DeviceID  string                 `json:"device_id"`
	Timestamp int64                  `json:"timestamp"`
	Metrics   map[string]interface{} `json:"metrics"`
}
