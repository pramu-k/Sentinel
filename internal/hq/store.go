package hq

import (
	"encoding/json"
	"time"
)

type Metric struct {
	Time       time.Time       `json:"time"`
	ServerID   string          `json:"server_id"`
	MetricType string          `json:"metric_type"`
	Resource   string          `json:"resource"`
	Value      float64         `json:"value"`
	Tags       json.RawMessage `json:"tags"`
}

type ServerStatus struct {
	ServerID  string    `json:"server_id"`
	LastSeen  time.Time `json:"last_seen"`
	IPAddress string    `json:"ip_address,omitempty"`
}

type ServiceStatus struct {
	ServiceName string    `json:"service_name"`
	Status      float64   `json:"status"`
	LastSeen    time.Time `json:"last_seen"`
}
