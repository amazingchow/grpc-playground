package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip"

	"github.com/amazingchow/dig-the-grpc/api"
)

// GrpcStreamServer gRPC流服务端
type GrpcStreamServer struct {
	logger zerolog.Logger
	server *grpc.Server
	l      net.Listener
	cfg    *GrpcStreamServerCfg
}

// GrpcStreamServerCfg gRPC流服务端配置
type GrpcStreamServerCfg struct {
	Port int    `json:"port"`
	Cert string `json:"cert"`
	Key  string `json:"key"`
}

// NewGrpcStreamServer 返回GrpcStreamServer实例.
func NewGrpcStreamServer(cfg *GrpcStreamServerCfg) (*GrpcStreamServer, error) {
	if cfg.Port == 0 {
		return nil, errors.Errorf("port must be specified")
	}

	srv := &GrpcStreamServer{}
	srv.logger = zerolog.New(os.Stdout).With().Str("from", "grpc stream server").Logger()
	srv.cfg = cfg

	return srv, nil
}

// Init 初始化gRPC流服务端.
func (srv *GrpcStreamServer) Init() error {
	var (
		opts  = []grpc.ServerOption{}
		creds credentials.TransportCredentials
		err   error
	)

	srv.l, err = net.Listen("tcp", fmt.Sprintf(":%d", srv.cfg.Port))
	if err != nil {
		return errors.Wrapf(err, "failed to listen on port %d", srv.cfg.Port)
	}

	if srv.cfg.Cert != "" && srv.cfg.Key != "" {
		creds, err = credentials.NewServerTLSFromFile(srv.cfg.Cert, srv.cfg.Key)
		if err != nil {
			return errors.Wrapf(err, "failed to create tls grpc server using cert '%s' and key '%s'", srv.cfg.Cert, srv.cfg.Key)
		}
		opts = append(opts, grpc.Creds(creds))
	}

	srv.server = grpc.NewServer(opts...)
	api.RegisterGrpcStreamServiceServer(srv.server, srv)

	return nil
}

// Run 开始运行gRPC流服务端.
func (srv *GrpcStreamServer) Run() {
	if err := srv.server.Serve(srv.l); err != nil {
		srv.logger.Error().Err(err)
	}
}

// Close 停止运行gRPC流服务端.
func (srv *GrpcStreamServer) Close() {
	if srv.server != nil {
		srv.server.GracefulStop()
	}
}

// Upload 实现文件传输接口.
func (srv *GrpcStreamServer) Upload(stream api.GrpcStreamService_UploadServer) error {
	var failed bool

RECV_LOOP:
	for {
		_, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				failed = false
			} else {
				srv.logger.Error().Err(err).Msg("failed unexpectedly while reading chunks from stream")
				failed = true
			}
			break RECV_LOOP
		}
	}

	if !failed {
		srv.logger.Info().Msg("upload successfully")

		if err := stream.SendAndClose(&api.UploadStatus{
			Message: "Successfully Upload",
			Code:    api.UploadStatusCode_STATUS_CODE_OK,
		}); err != nil {
			return errors.Wrapf(err, "failed to send status code")
		}
	} else {
		if err := stream.SendAndClose(&api.UploadStatus{
			Message: "Upload Failed",
			Code:    api.UploadStatusCode_STATUS_CODE_FAILED,
		}); err != nil {
			return errors.Wrapf(err, "failed to send status code")
		}
	}

	return nil
}
