package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
)

var (
	_Port = flag.Int("port", 8999, "server port")
	_Cert = flag.String("cert", "cert/cert.pem", "cert file")
	_Key  = flag.String("key", "cert/key.pem", "private key file")
)

func main() {
	flag.Parse()

	cfg := &GrpcStreamServerCfg{
		Port: *_Port,
		Cert: *_Cert,
		Key:  *_Key,
	}

	srv, err := NewGrpcStreamServer(cfg)
	if err != nil {
		panic(err)
	}
	if err = srv.Init(); err != nil {
		panic(err)
	}

	go srv.Run()
	defer func() {
		srv.Close()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
MAIN_LOOP:
	for { // nolint
		select {
		case <-sigCh:
			break MAIN_LOOP
		}
	}
}
