package grpc

import (
	"context"

	shippingpb "github.com/Laviniarm/microservices-proto/golang/shipping"
	"github.com/Laviniarm/microservices/shipping/internal/ports"
)

type Server struct {
	shippingpb.UnimplementedShippingServer
	app ports.ShippingService
}

func NewServer(app ports.ShippingService) *Server {
	return &Server{app: app}
}

func (s *Server) Estimate(ctx context.Context, req *shippingpb.EstimateRequest) (*shippingpb.EstimateResponse, error) {
	total := 0
	for _, it := range req.GetItems() {
		total += int(it.GetQuantity())
	}
	days, err := s.app.Estimate(ctx, total)
	if err != nil {
		return nil, err
	}
	return &shippingpb.EstimateResponse{EstimatedDays: days}, nil
}
