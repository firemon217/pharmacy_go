PROTO_DIR=proto
PB_DIR=pb

.PHONY: proto run

proto:
	mkdir -p $(PB_DIR)
	protoc --go_out=. --go_opt=module=app --go-grpc_out=. --go-grpc_opt=module=app $(PROTO_DIR)/pharmacy.proto

run: proto
	go run .
