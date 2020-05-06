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
	_Case = flag.Int("method", 1, `Test Cases:

1 - ONE-CONNECTION-PER-REQUEST
2 - ONLY-ONE-CONNECTION
3 - CONNECTION-POOL-WITH-EXPANSION`,
	)
	_Name = flag.String("name", "World", "name")
)

func main() {
	flag.Parse()

	var gRPCFunc func() string

	switch *_Case {
	case 1:
		{
			fmt.Println("using ONE-CONNECTION-PER-REQUEST test case...")

			gRPCFunc = func() string {
				conn, err := grpc.Dial("localhost:18081", grpc.WithInsecure())
				if err != nil {
					log.Fatalf("failed to connect: %v", err)
				}
				defer conn.Close()

				c := pb.NewGreeterClient(conn)

				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *_Name})
				if err != nil {
					log.Printf("failed to greet: %v", err)
					return "No Destination"
				}
				return r.Message
			}
		}
	case 2:
		{
			fmt.Println("using ONLY-ONE-CONNECTION test case...")

			conn, err := grpc.Dial("localhost:18081", grpc.WithInsecure())
			if err != nil {
				log.Fatalf("failed to connect: %v", err)
			}

			var clientPool = sync.Pool{
				New: func() interface{} {
					return pb.NewGreeterClient(conn)
				},
			}

			gRPCFunc = func() string {
				c := clientPool.Get().(pb.GreeterClient)
				defer clientPool.Put(c)

				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *_Name})
				if err != nil {
					log.Printf("failed to greet: %v", err)
					return "No Destination"
				}
				return r.Message
			}
		}
	case 3:
		{
			fmt.Println("using CONNECTION-POOL-WITH-EXPANSION test case...")

			var conns = make(chan *grpc.ClientConn, 5)
			for i := 0; i < 5; i++ {
				conn, err := grpc.Dial("localhost:18081", grpc.WithInsecure())
				if err != nil {
					log.Fatalf("failed to connect: %v", err)
				}
				conns <- conn
			}

			gRPCFunc = func() string {
				var conn *grpc.ClientConn
				var err error
				select {
				case conn = <-conns:
				default:
					conn, err = grpc.Dial("localhost:18081", grpc.WithInsecure())
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

				c := pb.NewGreeterClient(conn)

				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *_Name})
				if err != nil {
					log.Printf("failed to greet: %v", err)
					return "No Destination"
				}
				return r.Message
			}
		}
	default:
		{
			log.Fatalln("invalid input")
		}
	}

	http.HandleFunc("/performance", func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		rw.Write([]byte(gRPCFunc()))
	})

	fmt.Println("run :18080")
	http.ListenAndServe("localhost:18080", nil)
}
