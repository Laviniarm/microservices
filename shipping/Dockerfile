FROM golang:1.24.5 AS build
ENV GOTOOLCHAIN=auto
WORKDIR /src

# deps
COPY shipping/go.mod shipping/go.sum ./shipping/
RUN cd shipping && go mod download

# código
COPY shipping ./shipping

# build a partir da raiz do módulo
WORKDIR /src/shipping
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/shipping ./cmd/main.go

# runtime
FROM gcr.io/distroless/base-debian12
ENV SHIPPING_ADDR=":50053"
ENV ENV="production"
EXPOSE 50053
COPY --from=build /bin/shipping /shipping
ENTRYPOINT ["/shipping"]