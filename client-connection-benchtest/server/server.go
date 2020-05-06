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
	_Sleep = flag.Bool("sleep", false, "sleep for 200ms")
)

type myserver struct{}

func (s *myserver) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	if *_Sleep {
		time.Sleep(200 * time.Millisecond)
	}
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %s", req.Name)}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", "localhost:18081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &myserver{})
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
