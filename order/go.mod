module github.com/Laviniarm/microservices/order

go 1.24.0

require (
	github.com/Laviniarm/microservices-proto/golang/order v0.0.0-20250823012456-b25406c38caa
	google.golang.org/grpc v1.74.2
	gorm.io/driver/mysql v1.6.0
	gorm.io/gorm v1.30.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250811230008-5f3141c8851a // indirect
	google.golang.org/protobuf v1.36.7 // indirect
)

require (
	github.com/Laviniarm/microservices-proto/golang/payment v0.0.0-20250823012456-b25406c38caa
	github.com/Laviniarm/microservices-proto/golang/shipping v0.0.0-20250823012456-b25406c38caa
	github.com/go-sql-driver/mysql v1.9.3
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
)
