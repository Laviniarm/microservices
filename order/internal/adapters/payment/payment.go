package paymentadapter

import (
	"context"
	"log"
	"time"

	payment "github.com/Laviniarm/microservices-proto/golang/payment"
	"github.com/Laviniarm/microservices/order/internal/application/core/domain"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Adapter struct {
	payment payment.PaymentClient
}

func NewAdapter(paymentServiceURL string) (*Adapter, error) {
	// 🔄 Configuração do retry automático
	var opts []grpc.DialOption

	opts = append(opts,
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
			grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted), // apenas nesses casos
			grpc_retry.WithMax(5),                                            // até 5 tentativas
			grpc_retry.WithBackoff(grpc_retry.BackoffLinear(1*time.Second)),  // espera linear entre tentativas
		)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	conn, err := grpc.Dial(paymentServiceURL, opts...)
	if err != nil {
		return nil, err
	}

	client := payment.NewPaymentClient(conn)
	return &Adapter{payment: client}, nil
}

func (a *Adapter) Charge(order *domain.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := a.payment.Create(ctx, &payment.CreatePaymentRequest{
		UserId:     order.CustomerID,
		OrderId:    order.ID,
		TotalPrice: order.TotalPrice(),
	})

	if err != nil {
		// Verifica se foi erro de deadline
		if status.Code(err) == codes.DeadlineExceeded {
			log.Println("Erro: Timeout na comunicação com Payment (deadline excedido)")
		} else {
			log.Printf("Erro na chamada ao Payment: %v\n", err)
		}
		return err
	}

	return nil
}
