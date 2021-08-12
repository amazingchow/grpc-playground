package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/reflection"
)

var (
	sleepFlag = flag.Bool("sleep", false, "sleep for 200ms")
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (srv *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	if *sleepFlag {
		time.Sleep(200 * time.Millisecond)
	}
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %s", req.Name)}, nil
}

func main() {
	flag.Parse()

	l, err := net.Listen("tcp", "localhost:18889")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterGreeterServer(srv, &server{})
	reflection.Register(srv)

	if err = srv.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
