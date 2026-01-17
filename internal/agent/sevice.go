package agent

import (
	"context"
	"log"

	"github.com/kardianos/service"
)

type Program struct {
	Client *Client
	ctx    context.Context
	cancel context.CancelFunc
}

func NewProgram(client *Client) *Program {
	return &Program{
		Client: client,
	}
}

func (p *Program) Start(s service.Service) error {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	go p.run()
	return nil
}

func (p *Program) run() {
	defer p.cancel()
	// Start the gRPC client (this blocks until error or stop)
	if err := p.Client.Start(p.ctx); err != nil {
		log.Printf("Agent stopped with error: %v", err)
	}
}

func (p *Program) Stop(s service.Service) error {
	log.Println("Agent stopping...")
	p.cancel()
	p.Client.Stop()
	return nil
}
