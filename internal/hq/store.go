package hq

import (
	"context"
	"encoding/json"
	"sentinel/internal/proto"
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

type MetricStore interface {
	Init(ctx context.Context) error
	SaveBatch(ctx context.Context, batch *proto.MetricBatch, ipAddress string) error
	ListServers(ctx context.Context) ([]ServerStatus, error)
	GetMetrics(ctx context.Context, serverID string) ([]Metric, error)
	GetServiceStatus(ctx context.Context, serverID string) ([]ServiceStatus, error)
	Close()
}
