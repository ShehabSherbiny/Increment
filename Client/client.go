package main

import (
	"bufio"
	"context"
	service "mockexam/service"
	"log"
	"os"
	"fmt"
	"google.golang.org/grpc"
)

type Connection struct {
	port       string
	clientConn *grpc.ClientConn
	client     service.IncrementServiceClient
	context    context.Context
}

var connections []Connection

var serverPorts = []string{":9000", ":9001", ":9002"}

func main() {
	log.Printf("CLIENT STARTED")


	for i := range serverPorts {
		ctx, conn, client := setupConnection(i)
		newConn := Connection{
			port: serverPorts[i],
			clientConn: conn,
			client: client,
			context: ctx,
		}
		connections = append(connections, newConn)

		defer newConn.clientConn.Close()
	}
	
	for {
		promt()	
	}
}

func setupConnection(index int) (context.Context, *grpc.ClientConn, service.IncrementServiceClient) {
	conn, err := grpc.Dial(serverPorts[index], grpc.WithInsecure())

	if err != nil {
		log.Printf("Error: %v", err)
	}

	context := context.Background()

	client := service.NewIncrementServiceClient(conn)

	return context, conn, client
}

func promt() {
	log.Print("Write \"Increment\" to increment value")
	var counter = 0
	var shouldPrint = false
	var readContent string
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan() // Scanning user input
	readContent = scanner.Text()
	request := service.IncrementRequest{}
	ctx := context.Background()

	if readContent == "Increment" {
		fmt.Println()
		for i := range connections {
			if i == 0{
			value, err := connections[i].client.Increment(ctx, &request)
			if err != nil{
				connections[i] = connections[len(connections)-1]
				counter++
				shouldPrint = true
			}else{
				log.Printf("value: %v" , value.Value)
			}
			
		}else{
			value1, err := connections[i].client.Increment(ctx, &request)
			if err != nil{
				connections[i] = connections[len(connections)-1]
				counter++
			}else if shouldPrint{
				log.Printf("value: %v" , value1.Value)
				shouldPrint = false
			}
		}
	}
	connections = connections[:len(connections)-counter]
	}else{
		log.Println("Please write \"Increment\"")
	}
}