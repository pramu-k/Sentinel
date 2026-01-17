package agent

import (
	"context"
	"log"
	"sentinel/internal/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Config    *Config
	Collector *Collector
	Conn      *grpc.ClientConn
	Stream    proto.Sentinel_StreamMetricsClient
}

func (c *Client) Start(ctx context.Context) error {
	log.Printf("Connecting to HQ at %s...", c.Config.HQAddress)
	// For production, use credentials (TLS). For now, insecure is fine.
	conn, err := grpc.NewClient(c.Config.HQAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	c.Conn = conn
	client := proto.NewSentinelClient(conn)

	// Retry loop for stream connection
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := c.streamMetrics(ctx, client); err != nil {
				log.Printf("Stream error: %v. Retrying in 5s...", err)
				time.Sleep(5 * time.Second)
			}
		}
	}
}

func (c *Client) streamMetrics(ctx context.Context, client proto.SentinelClient) error {
	stream, err := client.StreamMetrics(ctx)
	if err != nil {
		return err
	}
	c.Stream = stream
	log.Println("Connected to HQ. Streaming metrics...")

	ticker := time.NewTicker(c.Config.CollectionInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return c.Stream.CloseSend()
		case <-ticker.C:
			batch := c.Collector.Collect()
			if err := c.Stream.Send(batch); err != nil {
				return err
			}
			log.Printf("Sent batch with %d metrics", len(batch.Metrics))
		}
	}
}

func (c *Client) Stop() {
	if c.Stream != nil {
		c.Stream.CloseSend()
	}
	if c.Conn != nil {
		c.Conn.Close()
	}
}
