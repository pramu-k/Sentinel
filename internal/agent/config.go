package agent

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type Config struct {
	HQAddress          string        `json:"hq_address"`
	CollectionInterval time.Duration `json:"-"`
	ServerID           string        `json:"server_id"`
	Services           []string      `json:"services"`
}

func LoadConfig() *Config {
	// Default config if config file is missing or invalid
	cfg := &Config{
		HQAddress:          "localhost:9090",
		CollectionInterval: 5 * time.Second,
		ServerID:           "winserv-01",
		Services:           []string{},
	}

	// Try to load from agent-config.json
	data, err := os.ReadFile("agent-config.json")
	if err == nil {
		type FileConfig struct {
			HQAddress          string   `json:"hq_address"`
			ServerID           string   `json:"server_id"`
			CollectionInterval string   `json:"collection_interval"`
			Services           []string `json:"services"`
		}
		var fCfg FileConfig
		if err := json.Unmarshal(data, &fCfg); err == nil {
			if fCfg.HQAddress != "" {
				cfg.HQAddress = fCfg.HQAddress
			}
			if fCfg.ServerID != "" {
				cfg.ServerID = fCfg.ServerID
			}
			if fCfg.CollectionInterval != "" {
				if d, err := time.ParseDuration(fCfg.CollectionInterval); err == nil {
					cfg.CollectionInterval = d
				} else {
					log.Printf("Invalid collection_interval '%s', using default 5s", fCfg.CollectionInterval)
				}
			}
			if len(fCfg.Services) > 0 {
				cfg.Services = fCfg.Services
			}
			log.Printf("Loaded config from agent-config.json: %+v", cfg)
		} else {
			log.Printf("Failed to parse agent-config.json: %v. Using defaults.", err)
		}
	} else {
		log.Println("agent-config.json not found. Using defaults.")
	}

	return cfg
}
