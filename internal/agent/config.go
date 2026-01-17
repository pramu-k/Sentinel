package agent

import "time"

type Config struct {
	HQAddress          string        `json:"hq_address"`
	CollectionInterval time.Duration `json:"-"`
	ServerID           string        `json:"server_id"`
	Services           []string      `json:"services"`
}
