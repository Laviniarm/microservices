package main

import (
	"log"

	"github.com/Laviniarm/microservices/order/config"
	"github.com/Laviniarm/microservices/order/internal/adapters/db"
	"github.com/Laviniarm/microservices/order/internal/adapters/grpc"
	paymentadapter "github.com/Laviniarm/microservices/order/internal/adapters/payment"
	"github.com/Laviniarm/microservices/order/internal/application/core/api"
)

func main() {
	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %v", err)
	}

	paymentAdapter, err := paymentadapter.NewAdapter(config.GetPaymentServiceUrl())
	if err != nil {
		log.Fatalf("Failed to initialize payment stub. Error: %v", err)
	}

	application := api.NewApplication(dbAdapter, paymentAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
