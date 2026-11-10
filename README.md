# Pharmacy microservice (Go kit)

Реализация варианта 10 «Аптека». HTTP и gRPC работают одновременно и используют один экземпляр сервиса и одно in-memory хранилище.

## Структура

- `implementation` — интерфейс сервиса и бизнес-логика.
- `transport/endpoints.go` — Go kit endpoints.
- `http` — HTTP/JSON transport.
- `transport/grpc.go` — Go kit gRPC transport и protobuf-адаптер.
- `repository` — потокобезопасное in-memory хранилище.
- `order` — доменная модель и ошибки.
- `proto` — контракт Protocol Buffers.
- `pb` — генерируемый Go-код (создаётся командой `make proto`).

## Требования

Установить Go 1.22+, `protoc`, плагины:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.5
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
```

Убедиться, что `$GOPATH/bin` добавлен в `PATH`, затем из корня проекта:

```bash
make proto
go mod tidy
go run .
```

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

Ответ содержит `id`. Использовать его в следующих запросах:

```bash
curl http://localhost:8080/medicines/ID
curl -i -X POST 'http://localhost:8080/medicines/ID/sell?count=3'
curl http://localhost:8080/medicines/ID
```

Альтернативно для продажи можно передать JSON `{"count":3}` в теле запроса.

## gRPC-проверка через grpcurl

Сервер reflection специально не включает: для проверки grpcurl использует исходный `.proto`:

```bash
grpcurl -plaintext -import-path . -proto proto/pharmacy.proto \
  -d '{"name":"Аспирин","manufacturer":"Bayer","price":120.5,"prescription":false,"quantity":20}' \
  localhost:9090 pharmacy.v1.Pharmacy/Create
```

Скопировать полученный `id`, затем:

```bash
grpcurl -plaintext -import-path . -proto proto/pharmacy.proto \
  -d '{"id":"ID"}' localhost:9090 pharmacy.v1.Pharmacy/GetByID

grpcurl -plaintext -import-path . -proto proto/pharmacy.proto \
  -d '{"id":"ID","count":3}' localhost:9090 pharmacy.v1.Pharmacy/Sell
```

## Коды ошибок

- HTTP `400` / gRPC `InvalidArgument` — некорректные данные или количество.
- HTTP `404` / gRPC `NotFound` — лекарство не найдено.
- HTTP `409` / gRPC `FailedPrecondition` — недостаточно товара на складе.
- HTTP `500` / gRPC `Internal` — внутренняя ошибка.

`Prescription` хранится как признак рецептурного препарата. Метод `Sell` по условию задания не получает сведения о наличии рецепта, поэтому проверка рецепта в нём не выполняется.
