# Pharmacy microservice (Go kit)

Реализация варианта 10 «Аптека». HTTP и gRPC работают одновременно и используют один экземпляр сервиса и одно in-memory хранилище.

Серверы:
- HTTP: `localhost:8080`
- gRPC: `localhost:9090`

## HTTP-проверка

Создать лекарство:

```bash
curl -i -X POST http://localhost:8080/medicines \
  -H 'Content-Type: application/json' \
  -d '{"name":"Аспирин","manufacturer":"Bayer","price":120.5,"prescription":false,"quantity":20}'
```

```bash
curl -i -X POST 'http://localhost:8080/medicines/ID/sell?count=3'
curl http://localhost:8080/medicines/ID
```

## gRPC-проверка через grpcurl

```bash
grpcurl -plaintext -import-path . -proto proto/pharmacy.proto \
  -d '{"name":"Аспирин","manufacturer":"Bayer","price":120.5,"prescription":false,"quantity":20}' \
  localhost:9090 pharmacy.v1.Pharmacy/Create
```

```bash
grpcurl -plaintext -import-path . -proto proto/pharmacy.proto \
  -d '{"id":"ID"}' localhost:9090 pharmacy.v1.Pharmacy/GetByID

grpcurl -plaintext -import-path . -proto proto/pharmacy.proto \
  -d '{"id":"ID","count":3}' localhost:9090 pharmacy.v1.Pharmacy/Sell
```
