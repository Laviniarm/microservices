package main

import (
	"log"

	"github.com/Laviniarm/microservices/order/config"
	"github.com/Laviniarm/microservices/order/internal/adapters/db"
	"github.com/Laviniarm/microservices/order/internal/adapters/grpc"
	"github.com/Laviniarm/microservices/order/internal/application/core/api"
)

func main() {
	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	application := api.NewApplication(dbAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
