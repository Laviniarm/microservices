package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" // driver para database/sql

	"github.com/Laviniarm/microservices/order/config"
	dbrepo "github.com/Laviniarm/microservices/order/internal/adapters/db"
	invrepo "github.com/Laviniarm/microservices/order/internal/adapters/db"
	ogrpc "github.com/Laviniarm/microservices/order/internal/adapters/grpc"
	paymentadapter "github.com/Laviniarm/microservices/order/internal/adapters/payment"
	paymentwrap "github.com/Laviniarm/microservices/order/internal/adapters/payment"
	shipcli "github.com/Laviniarm/microservices/order/internal/adapters/shipping"
	shipwrap "github.com/Laviniarm/microservices/order/internal/adapters/shippingwrap"
	"github.com/Laviniarm/microservices/order/internal/application/core/api"
)

func main() {
	dbAdapter, err := dbrepo.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %v", err)
	}

	sqlDB, err := sql.Open("mysql", config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("Failed to open sql DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping sql DB: %v", err)
	}
	inventory := invrepo.NewInventoryRepo(sqlDB)

	payAdp, err := paymentadapter.NewAdapter(config.GetPaymentServiceUrl())
	if err != nil {
		log.Fatalf("Failed to initialize payment client. Error: %v", err)
	}
	payment := paymentwrap.NewClientWrapper(payAdp)
	shipAdp, err := shipcli.NewAdapter()
	if err != nil {
		log.Fatal(err)
	}
	shipping := shipwrap.NewClientWrapper(shipAdp)

	application := api.NewApplication(dbAdapter, payAdp)

	grpcAdapter := ogrpc.NewAdapter(application, inventory, payment, shipping, config.GetApplicationPort())
	grpcAdapter.Run()
}
