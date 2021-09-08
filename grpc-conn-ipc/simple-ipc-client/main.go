package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

func main() {
	conn, err := grpc.Dial("unix:///var/run/grpc-conn-ipc-example.sock", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect, err: %v", err)
	}
	defer conn.Close()

	cli := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := cli.SayHello(ctx, &pb.HelloRequest{Name: "Grpc"})
	if err != nil {
		log.Fatalf("failed to greet, err: %v", err)
	}
	log.Println(resp.Message)
}
