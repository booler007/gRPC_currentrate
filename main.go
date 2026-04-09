package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/booler007/gRPC_currentrate/pb"
	"github.com/booler007/gRPC_currentrate/storage"
)

func main() {
	repo, err := storage.NewSQLiteRepository()
	if err != nil {
		log.Fatalf("storage init: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRatesServiceServer(grpcServer, NewRatesServer(repo))

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
