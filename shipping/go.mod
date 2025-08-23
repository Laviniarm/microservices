module github.com/Laviniarm/microservices/shipping

go 1.24.0

require (
    github.com/Laviniarm/microservices-proto/golang/shipping v0.0.0
    google.golang.org/grpc v1.75.0
    google.golang.org/protobuf v1.36.8
)

replace github.com/Laviniarm/microservices-proto/golang/shipping => ../../microservices-proto/golang/shipping
