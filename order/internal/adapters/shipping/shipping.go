package shipping

import (
	"context"
	"os"
	"time"

	shipping "github.com/Laviniarm/microservices-proto/golang/shipping"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	client shipping.ShippingClient
}

func NewAdapter() (*Adapter, error) {
	addr := os.Getenv("SHIPPING_ADDR")
	if addr == "" {
		addr = "shipping:50053"
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
			grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
			grpc_retry.WithMax(5),
			grpc_retry.WithBackoff(grpc_retry.BackoffLinear(1*time.Second)),
		)),
	}
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}

	return &Adapter{client: shipping.NewShippingClient(conn)}, nil
}

type Item struct {
	ID  string
	Qty int32
}

func (a *Adapter) Estimate(ctx context.Context, orderID int64, items []Item) (int32, error) {
	req := &shipping.EstimateRequest{
		OrderId: int32(orderID),
		Items:   make([]*shipping.Item, 0, len(items)),
	}
	for _, it := range items {
		req.Items = append(req.Items, &shipping.Item{Quantity: it.Qty})
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	resp, err := a.client.Estimate(ctx, req)
	if err != nil {
		return 0, err
	}
	return resp.GetEstimatedDays(), nil
}
