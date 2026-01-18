package hq

import (
	"io"
	"log"
	"sentinel/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type GRPCServer struct {
	proto.UnimplementedSentinelServer
	Store MetricStore
}

func NewGRPCServer(store MetricStore) *GRPCServer {
	return &GRPCServer{Store: store}
}

func (s *GRPCServer) Register(registrar grpc.ServiceRegistrar) {
	proto.RegisterSentinelServer(registrar, s)
}

func (s *GRPCServer) StreamMetrics(stream proto.Sentinel_StreamMetricsServer) error {
	ctx := stream.Context()
	ipAddress := "unknown"
	if p, ok := peer.FromContext(ctx); ok {
		ipAddress = p.Addr.String()
	}

	for {
		// Receive a batch
		batch, err := stream.Recv()
		if err == io.EOF {
			// Done reading
			return stream.SendAndClose(&proto.Ack{
				Success: true,
				Message: "Stream closed successfully",
			})
		}
		if err != nil {
			log.Printf("Error receiving stream: %v", err)
			return err
		}

		// Save to DB
		if err := s.Store.SaveBatch(ctx, batch, ipAddress); err != nil {
			log.Printf("Error saving batch from %s: %v", batch.ServerId, err)
		} else {
			log.Printf("Received & saved %d metrics from %s", len(batch.Metrics), batch.ServerId)
		}
	}
}
