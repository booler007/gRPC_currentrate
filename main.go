package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("gRPC server listening on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("serve: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down...")

	grpcServer.GracefulStop()

	if err := repo.Close(); err != nil {
		log.Printf("db close: %v", err)
	}

	log.Println("shutdown complete")
}
