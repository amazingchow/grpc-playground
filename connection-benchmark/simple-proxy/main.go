package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

var (
	caseFlag = flag.Int("method", 1, `Test Cases:

1) ONE_CONNECTION_PER_REQUEST
2) ONLY_ONE_CONNECTION
3) CONNECTION_POOL_WITH_EXPANSION`,
	)
)

const (
	ONE_CONNECTION_PER_REQUEST     = 1
	ONLY_ONE_CONNECTION            = 2
	CONNECTION_POOL_WITH_EXPANSION = 3
)

func main() {
	flag.Parse()

	var gRPCHandler func() string

	switch *caseFlag {
	case ONE_CONNECTION_PER_REQUEST:
		{
			fmt.Println("using ONE_CONNECTION_PER_REQUEST test case...")

			gRPCHandler = func() string {
				conn, err := grpc.Dial("localhost:18889", grpc.WithInsecure())
				if err != nil {
					log.Fatalf("failed to connect: %v", err)
				}
				defer conn.Close()

				cli := pb.NewGreeterClient(conn)

				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				resp, err := cli.SayHello(ctx, &pb.HelloRequest{Name: "Grpc"})
				if err != nil {
					log.Printf("failed to greet: %v", err)
					return "No Destination"
				}
				return resp.Message
			}
		}
	case ONLY_ONE_CONNECTION:
		{
			fmt.Println("using ONLY_ONE_CONNECTION test case...")

			conn, err := grpc.Dial("localhost:18889", grpc.WithInsecure())
			if err != nil {
				log.Fatalf("failed to connect: %v", err)
			}

			var clientPool = sync.Pool{
				New: func() interface{} {
					return pb.NewGreeterClient(conn)
				},
			}

			gRPCHandler = func() string {
				cli := clientPool.Get().(pb.GreeterClient)
				defer clientPool.Put(cli)

				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				resp, err := cli.SayHello(ctx, &pb.HelloRequest{Name: "Grpc"})
				if err != nil {
					log.Printf("failed to greet: %v", err)
					return "No Destination"
				}
				return resp.Message
			}
		}
	case CONNECTION_POOL_WITH_EXPANSION:
		{
			fmt.Println("using CONNECTION_POOL_WITH_EXPANSION test case...")

			var conns = make(chan *grpc.ClientConn, 5)
			for i := 0; i < 5; i++ {
				conn, err := grpc.Dial("localhost:18889", grpc.WithInsecure())
				if err != nil {
					log.Fatalf("failed to connect: %v", err)
				}
				conns <- conn
			}

			gRPCHandler = func() string {
				var conn *grpc.ClientConn
				var err error

				select {
				case conn = <-conns:
				default:
					conn, err = grpc.Dial("localhost:18889", grpc.WithInsecure())
					if err != nil {
						log.Fatalf("failed to connect: %v", err)
					}
				}
				defer func() {
					select {
					case conns <- conn:
					default:
						conn.Close()
					}
				}()

				cli := pb.NewGreeterClient(conn)

				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				resp, err := cli.SayHello(ctx, &pb.HelloRequest{Name: "Grpc"})
				if err != nil {
					log.Printf("failed to greet: %v", err)
					return "No Destination"
				}
				return resp.Message
			}
		}
	default:
		{
			log.Fatalln("invalid input")
		}
	}

	http.HandleFunc("/performance", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(gRPCHandler())) // nolint
	})
	http.ListenAndServe("localhost:18888", nil) // nolint
}
