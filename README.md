# Microsserviços com gRPC (Order • Payment • Shipping)

Implementação da parte final da prática: microsserviço **Shipping** com gRPC, integração com **Order** (apenas após pagamento aprovado), validação de **estoque** no Order via MySQL e **deploy** com Docker.

## Arquitetura

- **order**: recebe o pedido, valida itens contra o estoque (MySQL), persiste pedido, chama **payment**; se aprovado, chama **shipping** para estimar prazo.
- **payment**: serviço gRPC dummy que aprova o pagamento (pode ser adaptado).
- **shipping**: calcula prazo de entrega a partir do **total de unidades** (mínimo 1 dia, e **+1 dia a cada 5 unidades**).

## Requisitos

- Go 1.24+ (nas imagens Docker)
- Docker e Docker Compose
- `grpcurl` instalado no host para testes

---

## Como subir com Docker

> Rode tudo a partir da **raiz** deste repositório (onde está o `docker-compose.yml`).

(Opcional, mas recomendado em 1ª execução) limpar volumes:
```bash
docker compose down -v
```
1) Subir containers:
```bash
docker compose up -d --build
```

Banco de Dados (MySQL) - Verificar rapidamente:
```bash
docker exec -it microservices-mysql-1 mysql -uroot -proot -e \
"USE orders; SHOW TABLES; SELECT * FROM inventory_items;"
```

2) Teste do Shipping (isolado)
```bash
grpcurl -plaintext -d '{
  "order_id": 123,
  "items": [{"quantity": 6}, {"quantity": 4}]
}' localhost:50053 Shipping/Estimate
```
3) Teste Order – validação de estoque (SKU inexistente → erro apropriado)
```bash
grpcurl -plaintext -d '{
  "costumerId": 1,
  "orderItems": [
    {"productCode":"SKU-404","unitPrice":10.0,"quantity":1}
  ],
  "totalPrice": 10.0
}' localhost:50051 Order/Create
```
4) Teste Order – fluxo Ok (pagamento aprovado → chama Shipping)
```bash
grpcurl -plaintext -d '{
  "costumerId": 1,
  "orderItems": [
    {"productCode":"SKU-1","unitPrice":50.0,"quantity":6}
  ],
  "totalPrice": 300.0
}' localhost:50051 Order/Create
```

### 4) Provar que o Order só chama Shipping se Payment aprovar

#### Passo 1: Simular falha de pagamento  
Pare o serviço **payment**:  
```bash
docker compose stop payment

grpcurl -plaintext -d '{
  "costumerId": 1,
  "orderItems": [{"productCode":"SKU-1","unitPrice":50.0,"quantity":2}],
  "totalPrice": 100.0
}' localhost:50051 Order/Create
```
