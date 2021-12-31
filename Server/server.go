package main

import (
	"context"
	"log"
	service "mockexam/service"
	"mockexam/utils"
	"net"
	"os"

	"google.golang.org/grpc"
)


type Server struct {
	service.UnimplementedIncrementServiceServer
	counter utils.Counter  
}

// Example of run: go run . 9000
func main() {
	log.Printf("SERVER STARTED")

	args := os.Args

	port := ":" + args[1]

	go setupServer(port)

	for {}
}

func setupServer(port string) {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	grpcServer := grpc.NewServer()
	s := Server{}
	
	service.RegisterIncrementServiceServer(grpcServer, &s)

	log.Printf("Listening at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *Server) Increment(ctx context.Context, request *service.IncrementRequest) (*service.ValueReturn, error) {
	value := s.counter.Value()
	log.Printf("Incrementing value to... %v", value+1)
	s.counter.Increment()
	return &service.ValueReturn{Value: int32(value)} ,nil
}