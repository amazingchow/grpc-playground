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

	"github.com/amazingchow/grpc-playground/grpc-file-transfer-tool/api"
)

// GrpcStreamServer gRPC流服务端
type GrpcStreamServer struct {
	logger zerolog.Logger
	cfg    *GrpcStreamServerCfg
	srv    *grpc.Server
	l      net.Listener
}

// GrpcStreamServerCfg gRPC流服务端配置
type GrpcStreamServerCfg struct {
	Port int    `json:"port"`
	Cert string `json:"cert"`
	Key  string `json:"key"`
}

// NewGrpcStreamServer 返回GrpcStreamServer实例.
func NewGrpcStreamServer(cfg *GrpcStreamServerCfg) (*GrpcStreamServer, error) {
	srv := &GrpcStreamServer{}
	srv.logger = zerolog.New(os.Stdout).With().Str("from", "grpc stream server").Logger()
	srv.cfg = cfg
	return srv, nil
}

// Init 初始化gRPC流服务端.
func (gsrv *GrpcStreamServer) Init() error {
	var (
		opts = []grpc.ServerOption{}
		err  error
	)

	gsrv.l, err = net.Listen("tcp", fmt.Sprintf(":%d", gsrv.cfg.Port))
	if err != nil {
		gsrv.logger.Error().Err(err).Msgf("failed to listen on port %d", gsrv.cfg.Port)
		return errors.Wrapf(err, "failed to listen on port %d", gsrv.cfg.Port)
	}

	if gsrv.cfg.Cert != "" && gsrv.cfg.Key != "" {
		creds, err := credentials.NewServerTLSFromFile(gsrv.cfg.Cert, gsrv.cfg.Key)
		if err != nil {
			gsrv.logger.Error().Err(err).Msgf("failed to create tls-grpc-server using cert '%s' and key '%s'", gsrv.cfg.Cert, gsrv.cfg.Key)
			return errors.Wrapf(err, "failed to create tls-grpc-server using cert '%s' and key '%s'", gsrv.cfg.Cert, gsrv.cfg.Key)
		}
		opts = append(opts, grpc.Creds(creds))
	}

	gsrv.srv = grpc.NewServer(opts...)
	api.RegisterGrpcStreamServiceServer(gsrv.srv, gsrv)

	return nil
}

// Run 开始运行gRPC流服务端.
func (gsrv *GrpcStreamServer) Run() {
	if err := gsrv.srv.Serve(gsrv.l); err != nil {
		gsrv.logger.Error().Err(err)
	}
}

// Close 停止运行gRPC流服务端.
func (gsrv *GrpcStreamServer) Close() {
	if gsrv.srv != nil {
		gsrv.srv.GracefulStop()
	}
}

// Upload 实现文件传输接口.
func (gsrv *GrpcStreamServer) Upload(stream api.GrpcStreamService_UploadServer) error {
	var failed bool

RECV_LOOP:
	for {
		_, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				failed = false
			} else {
				gsrv.logger.Error().Err(err).Msg("failed unexpectedly while reading chunks from stream")
				failed = true
			}
			break RECV_LOOP
		}
		// TODO: store the uploaded file
	}

	if !failed {
		gsrv.logger.Info().Msg("upload successfully")

		if err := stream.SendAndClose(&api.UploadStatus{
			Message: "Successfully Upload",
			Code:    api.UploadStatusCode_STATUS_CODE_OK,
		}); err != nil {
			gsrv.logger.Error().Err(err).Msg("failed to send status code")
			return errors.Wrapf(err, "failed to send status code")
		}
	} else {
		if err := stream.SendAndClose(&api.UploadStatus{
			Message: "Upload Failed",
			Code:    api.UploadStatusCode_STATUS_CODE_FAILED,
		}); err != nil {
			gsrv.logger.Error().Err(err).Msg("failed to send status code")
			return errors.Wrapf(err, "failed to send status code")
		}
	}

	return nil
}
