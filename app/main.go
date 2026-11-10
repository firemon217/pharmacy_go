package main

import (
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"

	httptransport "pharmacy/app/http"
	"pharmacy/app/implementation"
	"pharmacy/app/pb"
	"pharmacy/app/repository"
	"pharmacy/app/transport"
)

func main() {
	repo := repository.NewMemory()
	svc := implementation.NewService(repo)
	endpoints := transport.MakeEndpoints(svc)

	// HTTP transport.
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: httptransport.NewHTTPHandler(endpoints),
	}
	go func() {
		log.Println("HTTP server listening on :8080")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server: %v", err)
		}
	}()

	// gRPC transport.
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("gRPC listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPharmacyServer(grpcServer, transport.NewGRPCServer(endpoints))

	log.Println("gRPC server listening on :9090")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("gRPC server: %v", err)
	}
}
