package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"

	pb "github.com/Laviniarm/microservices-proto/golang/order"
	"github.com/Laviniarm/microservices/order/config"
	"github.com/Laviniarm/microservices/order/internal/application/core/domain"
	"github.com/Laviniarm/microservices/order/internal/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Adapter struct {
	api  ports.APIPort
	port int
	pb.UnimplementedOrderServer
}

func NewAdapter(api ports.APIPort, port int) *Adapter {
	return &Adapter{api: api, port: port}
}

func (a *Adapter) Create(ctx context.Context, request *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	var totalItems int32 = 0
	for _, orderItem := range request.OrderItems {
		totalItems += orderItem.Quantity
	}
	
	if totalItems > 50 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"pedido não pode ser realizado: quantidade total de itens (%d) excede o limite de 50", totalItems,
		)
	}
	var orderItems []domain.OrderItem
	for _, orderItem := range request.OrderItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}

	newOrder := domain.NewOrder(int64(request.CustomerId), orderItems)
	result, err := a.api.PlaceOrder(newOrder)
	if err != nil {
		return nil, err
	}

	return &pb.CreateOrderResponse{OrderId: int32(result.ID)}, nil
}

func (a *Adapter) Run() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatalf("failed to listen on port %d: %v", a.port, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServer(grpcServer, a)

	if config.GetEnv() == "development" {
		reflection.Register(grpcServer)
	}

	log.Printf("gRPC server running on port %d", a.port)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
	}
}
