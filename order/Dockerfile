FROM golang:1.24.5 AS build
ENV GOTOOLCHAIN=auto
WORKDIR /src

# deps (usa o go.mod do módulo order)
COPY order/go.mod order/go.sum ./order/
RUN cd order && go mod download

# código
COPY order ./order

# build a partir da raiz do módulo (onde está o go.mod)
WORKDIR /src/order
# se seu main é order/cmd/main.go:
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/order ./cmd/main.go
# (alternativa equivalente)
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/order ./cmd

# runtime
FROM gcr.io/distroless/base-debian12
ENV ENV="production" \
    DATA_SOURCE_URL="root:root@tcp(mysql:3306)/orders?parseTime=true" \
    PAYMENT_SERVICE_URL="payment:50052" \
    SHIPPING_ADDR="shipping:50053" \
    APPLICATION_PORT="50051"
EXPOSE 50051
COPY --from=build /bin/order /order
ENTRYPOINT ["/order"]