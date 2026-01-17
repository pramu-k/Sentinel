package main

import (
	"log"
	"os"

	"sentinel/internal/agent"

	"github.com/kardianos/service"
)

func main() {
	// 1. Setup Service Configuration
	svcConfig := &service.Config{
		Name:        "SentinelAgent",
		DisplayName: "Sentinel Monitoring Agent",
		Description: "Collects and streams system metrics to HQ.",
	}

	// 2. Initialize Internal Components
	cfg := agent.LoadConfig()
	collector := agent.NewCollector(cfg)
	client := agent.NewClient(cfg, collector)
	prg := agent.NewProgram(client)

	// 3. Create Service
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	// 4. Handle Control Actions (install/uninstall/start/stop)
	if len(os.Args) > 1 {
		err = service.Control(s, os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// 5. Run Service
	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}

}
