package agent

import (
	"log"
	"math"
	"strings"

	"sentinel/internal/proto"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
	"google.golang.org/protobuf/types/known/timestamppb"
)
type Collector struct {
	Config *Config
}
func NewCollector(cfg *Config) *Collector {
	return &Collector{Config: cfg}
}

func (c *Collector) Collect() *proto.MetricBatch {
	batch := &proto.MetricBatch{
		ServerId:  c.Config.ServerID,
		Timestamp: timestamppb.Now(),
		Metrics:   []*proto.Metric{},
	}

	// 1. CPU Usage
	percent, err := cpu.Percent(0, false)
	if err == nil && len(percent) > 0 {
		batch.Metrics = append(batch.Metrics, &proto.Metric{
			Type:  "cpu_usage",
			Value: round(percent[0]),
			Tags:  map[string]string{"unit": "percent"},
		})
	} else {
		log.Printf("Error getting CPU: %v", err)
	}

	// 2. Memory Usage
	v, err := mem.VirtualMemory()
	if err == nil {
		batch.Metrics = append(batch.Metrics, &proto.Metric{
			Type:  "memory_used_percent",
			Value: round(v.UsedPercent),
			Tags:  map[string]string{"unit": "percent"},
		})
		batch.Metrics = append(batch.Metrics, &proto.Metric{
			Type:  "memory_free_mb",
			Value: float64(v.Free) / 1024 / 1024,
			Tags:  map[string]string{"unit": "mb"},
		})
		batch.Metrics = append(batch.Metrics, &proto.Metric{
			Type:  "memory_total_mb",
			Value: float64(v.Total) / 1024 / 1024,
			Tags:  map[string]string{"unit": "mb"},
		})
		batch.Metrics = append(batch.Metrics, &proto.Metric{
			Type:  "memory_available_mb",
			Value: float64(v.Available) / 1024 / 1024,
			Tags:  map[string]string{"unit": "mb"},
		})
	} else {
		log.Printf("Error getting Memory: %v", err)
	}

	// 3. Disk Usage (Root)
	d, err := disk.Usage("/")
	if err == nil {
		batch.Metrics = append(batch.Metrics, &proto.Metric{
			Type:  "disk_free_gb",
			Value: float64(d.Free) / 1024 / 1024 / 1024,
			Tags:  map[string]string{"unit": "gb", "path": "/"},
		})
		batch.Metrics = append(batch.Metrics, &proto.Metric{
			Type:  "disk_used_percent",
			Value: round(d.UsedPercent),
			Tags:  map[string]string{"unit": "percent", "path": "/"},
		})
	} else {
		log.Printf("Error getting Disk: %v", err)
	}

	// 4. Service Monitoring
	if len(c.Config.Services) > 0 {
		procs, err := process.Processes()
		if err == nil {
			// Create a set of services for fast lookup
			servicesToMonitor := make(map[string]bool)
			for _, s := range c.Config.Services {
				servicesToMonitor[strings.ToLower(s)] = true
			}

			for _, p := range procs {
				name, err := p.Name()
				if err != nil {
					continue
				}
				nameLower := strings.ToLower(name)

				// Simple contains check (e.g. "postgres.exe" contains "postgres")
				foundService := ""
				for s := range servicesToMonitor {
					if strings.Contains(nameLower, s) {
						foundService = s
						break
					}
				}

				if foundService != "" {
					// Found a monitored service!
					cpuPercent, _ := p.CPUPercent()
					memInfo, _ := p.MemoryInfo()

					batch.Metrics = append(batch.Metrics, &proto.Metric{
						Type:  "service_cpu",
						Value: round(cpuPercent),
						Tags:  map[string]string{"service": foundService, "pid": string(rune(p.Pid))},
					})

					if memInfo != nil {
						batch.Metrics = append(batch.Metrics, &proto.Metric{
							Type:  "service_memory_mb",
							Value: float64(memInfo.RSS) / 1024 / 1024,
							Tags:  map[string]string{"service": foundService},
						})
					}

					// We could also report "up" status implicitly by the presence of these metrics
					batch.Metrics = append(batch.Metrics, &proto.Metric{
						Type:  "service_status",
						Value: 1, // 1 = Up
						Tags:  map[string]string{"service": foundService},
					})
				}
			}
		}
	}

	return batch
}

func round(val float64) float64 {
	return math.Round(val*100) / 100
}