package cmd

import (
	"github.com/Laviniarm/microservices/shipping/internal/application/core/api"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	shippingpb "github.com/Laviniarm/microservices-proto/golang/shipping"
	g "github.com/Laviniarm/microservices/shipping/internal/adapters/grpc"
	_ "github.com/Laviniarm/microservices/shipping/internal/application/core/api"
)

func main() {
	addr := getEnv("SHIPPING_ADDR", ":50053")

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	app := api.New()
	s := grpc.NewServer()
	shippingpb.RegisterShippingServer(s, g.NewServer(app))
	// habilita reflection pra facilitar o grpcurl em dev
	if os.Getenv("ENV") == "development" {
		reflection.Register(s)
	}

	log.Printf("shipping listening at %s", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

func getEnv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
