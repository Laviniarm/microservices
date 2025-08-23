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

// --- Novas portas (interfaces) esperadas ---
// Ajuste os imports se você já tiver clients concretos prontos em packages próprios.
//
// PaymentClient deve encapsular a chamada gRPC para o microsserviço Payment.
// Retorna approved=true/false e um erro (erro de infra, timeout, etc).
type PaymentClient interface {
	Process(ctx context.Context, req struct {
		OrderID    int64
		CustomerID int64
		Items      []domain.OrderItem
		Total      float32
		Currency   string
	}) (approved bool, err error)
}

// ShippingClient deve encapsular a chamada gRPC para o microsserviço Shipping.
// Retorna o número de dias estimados para entrega.
type ShippingClient interface {
	Estimate(ctx context.Context, orderID int64, items []struct {
		ID  string
		Qty int32
	}) (estimatedDays int32, err error)
}

// Adapter agora recebe as novas dependências:
// - api: sua aplicação (PlaceOrder persiste a ordem)
// - inventory: repositório para validar SKUs (existência em DB)
// - payment: client gRPC do Payment
// - shipping: client gRPC do Shipping
type Adapter struct {
	api       ports.APIPort
	inventory ports.InventoryRepository
	payment   PaymentClient
	shipping  ShippingClient

	port int
	pb.UnimplementedOrderServer
}

// NewAdapter atualizado para injetar as dependências.
func NewAdapter(
	api ports.APIPort,
	inventory ports.InventoryRepository,
	payment PaymentClient,
	shipping ShippingClient,
	port int,
) *Adapter {
	return &Adapter{
		api:       api,
		inventory: inventory,
		payment:   payment,
		shipping:  shipping,
		port:      port,
	}
}

func (a *Adapter) Create(ctx context.Context, request *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	// 0) Regra já existente: limitar quantidade total
	var totalItems int32
	for _, it := range request.OrderItems {
		totalItems += it.Quantity
	}
	if totalItems > 50 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"pedido não pode ser realizado: quantidade total de itens (%d) excede o limite de 50", totalItems,
		)
	}

	// 1) Validar SKUs no estoque (requisito adicional - Order)
	ids := make([]string, 0, len(request.OrderItems))
	for _, it := range request.OrderItems {
		ids = append(ids, it.ProductCode)
	}
	missing, err := a.inventory.MissingIDs(ctx, ids)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "inventory check failed: %v", err)
	}
	if len(missing) > 0 {
		return nil, status.Errorf(codes.InvalidArgument, "unknown item IDs: %v", missing)
	}

	// 2) Montar itens de domínio
	orderItems := make([]domain.OrderItem, 0, len(request.OrderItems))
	for _, it := range request.OrderItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: it.ProductCode,
			UnitPrice:   it.UnitPrice,
			Quantity:    it.Quantity,
		})
	}

	// 3) Persistir/registrar a ordem (PlaceOrder salva e devolve ID)
	newOrder := domain.NewOrder(int64(request.CustomerId), orderItems)
	result, err := a.api.PlaceOrder(newOrder)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to persist order: %v", err)
	}

	// 4) Processar pagamento
	//    Ajuste o payload conforme seu Payment.proto/método real.
	payApproved, err := a.payment.Process(ctx, struct {
		OrderID    int64
		CustomerID int64
		Items      []domain.OrderItem
		Total      float32
		Currency   string
	}{
		OrderID:    result.ID,
		CustomerID: int64(request.CustomerId),
		Items:      orderItems,
		Total:      request.TotalPrice, // ou recalcule somando UnitPrice*Quantity
		Currency:   "BRL",
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "payment error: %v", err)
	}
	if !payApproved {
		return nil, status.Error(codes.FailedPrecondition, "payment not approved")
	}

	// 5) Chamar Shipping APENAS se o pagamento for aprovado
	itemsForShipping := make([]struct {
		ID  string
		Qty int32
	}, 0, len(request.OrderItems))
	for _, it := range request.OrderItems {
		itemsForShipping = append(itemsForShipping, struct {
			ID  string
			Qty int32
		}{
			ID:  it.ProductCode,
			Qty: it.Quantity,
		})
	}

	estimatedDays, err := a.shipping.Estimate(ctx, result.ID, itemsForShipping)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "shipping estimate failed: %v", err)
	}

	// 6) Montar resposta
	resp := &pb.CreateOrderResponse{
		OrderId: int32(result.ID),
	}

	// Se você adicionou o campo no seu order.proto:
	//   message CreateOrderResponse {
	//     int32 order_id = 1;
	//     int32 estimated_delivery_days = 2;
	//   }
	// então descomente a linha abaixo:
	//
	// resp.EstimatedDeliveryDays = estimatedDays

	// Se o campo ainda não existe no proto, você pode só logar o valor por ora:
	_ = estimatedDays // remover quando adicionar ao proto

	return resp, nil
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
