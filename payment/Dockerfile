FROM golang:1.24.5 AS build
ENV GOTOOLCHAIN=auto
WORKDIR /src

# deps (usa o go.mod do módulo payment)
COPY payment/go.mod payment/go.sum ./payment/
RUN cd payment && go mod download

# código
COPY payment ./payment

# build a partir da raiz do módulo
WORKDIR /src/payment
# se seu main é payment/cmd/main.go:
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/payment ./cmd/main.go
# (alternativa equivalente)
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/payment ./cmd

# runtime
FROM gcr.io/distroless/base-debian12
ENV PAYMENT_ADDR=":50052"
EXPOSE 50052
COPY --from=build /bin/payment /payment
ENTRYPOINT ["/payment"]