package main

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

func UnixConnect(context.Context, string) (net.Conn, error) {
	addr, err := net.ResolveUnixAddr("unix", "/run/grpc-conn-ipc-example.sock")
	if err != nil {
		log.Printf("failed to resolve unix addr, err: %v\n", err)
		return nil, err
	}
	conn, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		log.Printf("failed to dial unix, err: %v\n", err)
		return nil, err
	}
	return conn, nil
}

func main() {
	conn, err := grpc.Dial("/run/grpc-conn-ipc-example.sock", grpc.WithInsecure(), grpc.WithContextDialer(UnixConnect))
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
