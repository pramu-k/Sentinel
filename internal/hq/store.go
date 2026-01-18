package hq

import (
	"context"
	"encoding/json"
	"fmt"
	"sentinel/internal/proto"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
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

type DBStore struct {
	db *pgxpool.Pool
}

func NewDBStore(ctx context.Context, connString string) (*DBStore, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	return &DBStore{db: pool}, nil
}

func (s *DBStore) Close() {
	s.db.Close()
}

func (s *DBStore) Init(ctx context.Context) error {
	// 1. Create Metrics Table (New Schema with Resource)
	_, err := s.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS metrics (
			time        TIMESTAMPTZ NOT NULL,
			server_id   TEXT NOT NULL,
			metric_type TEXT NOT NULL,
			resource    TEXT NOT NULL DEFAULT '',
			value       DOUBLE PRECISION NOT NULL,
			tags        JSONB,
			CONSTRAINT metrics_pkey PRIMARY KEY (server_id, metric_type, resource, time)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create metrics table: %w", err)
	}

	// 2. Create Server Status Table
	_, err = s.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS server_status (
			server_id   TEXT PRIMARY KEY,
			last_seen   TIMESTAMPTZ NOT NULL,
			ip_address  TEXT
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create server_status table: %w", err)
	}

	return nil
}

func (s *DBStore) SaveBatch(ctx context.Context, batch *proto.MetricBatch, ipAddress string) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Update Last Seen
	_, err = tx.Exec(ctx, `
		INSERT INTO server_status (server_id, last_seen, ip_address)
		VALUES ($1, $2, $3)
		ON CONFLICT (server_id) DO UPDATE SET last_seen = $2, ip_address = $3
	`, batch.ServerId, batch.Timestamp.AsTime(), ipAddress)
	if err != nil {
		return err
	}

	// Insert Metrics
	for _, m := range batch.Metrics {
		tagsJSON, _ := json.Marshal(m.Tags)

		// Extract Resource from tags (for uniqueness)
		resource := ""
		if val, ok := m.Tags["service"]; ok {
			resource = val
		} else if val, ok := m.Tags["path"]; ok {
			resource = val
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO metrics (time, server_id, metric_type, resource, value, tags)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (server_id, metric_type, resource, time) DO NOTHING
		`, batch.Timestamp.AsTime(), batch.ServerId, m.Type, resource, m.Value, tagsJSON)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *DBStore) ListServers(ctx context.Context) ([]ServerStatus, error) {
	rows, err := s.db.Query(ctx, "SELECT server_id, last_seen, COALESCE(ip_address, '') FROM server_status ORDER BY last_seen DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []ServerStatus
	for rows.Next() {
		var s ServerStatus
		if err := rows.Scan(&s.ServerID, &s.LastSeen, &s.IPAddress); err != nil {
			return nil, err
		}
		servers = append(servers, s)
	}
	return servers, nil
}

func (s *DBStore) GetMetrics(ctx context.Context, serverID string) ([]Metric, error) {
	// Get last 100 metrics for this server
	rows, err := s.db.Query(ctx, `
		SELECT time, server_id, metric_type, resource, value, tags 
		FROM metrics 
		WHERE server_id = $1 
		ORDER BY time DESC 
		LIMIT 100
	`, serverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []Metric
	for rows.Next() {
		var m Metric
		if err := rows.Scan(&m.Time, &m.ServerID, &m.MetricType, &m.Resource, &m.Value, &m.Tags); err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}

func (s *DBStore) GetServiceStatus(ctx context.Context, serverID string) ([]ServiceStatus, error) {
	rows, err := s.db.Query(ctx, `
		SELECT DISTINCT ON (resource) 
			resource as service_name,
			value as status,
			time as last_seen
		FROM metrics
		WHERE server_id = $1 
			AND metric_type = 'service_status'
			AND resource != ''
		ORDER BY resource, time DESC
	`, serverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []ServiceStatus
	for rows.Next() {
		var s ServiceStatus
		if err := rows.Scan(&s.ServiceName, &s.Status, &s.LastSeen); err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, nil
}
