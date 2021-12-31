package main

import (
	"context"
	"log"
	service "mockexam/service"
	"net"
	"os"

	"google.golang.org/grpc"
)


type Server struct {
	service.UnimplementedIncrementServiceServer
	value  chan int32
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
	s.value = make(chan int32, 1)
	s.value <- 0
	service.RegisterIncrementServiceServer(grpcServer, &s)

	log.Printf("Listening at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *Server) Increment(ctx context.Context, request *service.IncrementRequest) (*service.ValueReturn, error) {
	value  := <- s.value
	log.Printf("Incrementing value to... %v", value+1)
	s.value <- value+1
	log.Printf("%v", value)
	return &service.ValueReturn{Value: value} ,nil
}