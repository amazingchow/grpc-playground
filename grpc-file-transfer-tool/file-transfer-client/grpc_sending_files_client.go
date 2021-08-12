package main

import (
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip"

	"github.com/amazingchow/grpc-playground/grpc-file-transfer-tool/api"
	"github.com/amazingchow/grpc-playground/grpc-file-transfer-tool/common"
)

// GRPCStreamClient gRPC流客户端
type GRPCStreamClient struct {
	logger zerolog.Logger
	cfg    *GRPCStreamClientCfg
	client api.GrpcStreamServiceClient
	conn   *grpc.ClientConn
}

// GRPCStreamClientCfg gRPC流客户端配置
type GRPCStreamClientCfg struct {
	Address    string `json:"address"`
	ChunkSize  int    `json:"chunk_size"`
	Compressed bool   `json:"compressed"`
	RootCert   string `json:"root_cert"`
}

// NewGRPCStreamClient 返回GRPCStreamClient实例.
func NewGRPCStreamClient(cfg *GRPCStreamClientCfg) (*GRPCStreamClient, error) {
	var (
		opts = []grpc.DialOption{}
		err  error
	)

	if cfg.Address == "" {
		return nil, errors.Errorf("address must be specified")
	}
	if cfg.ChunkSize <= 0 {
		return nil, errors.Errorf("chunk_size must be specified")
	} else if cfg.ChunkSize > (1 << 22) {
		return nil, errors.Errorf("chunk_size must be less than 4MB")
	}
	if cfg.Compressed {
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")))
	}
	if cfg.RootCert != "" {
		creds, err := credentials.NewClientTLSFromFile(cfg.RootCert, "SummyChou") // change the serverNameOverride for yourself
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create tls-grpc-client using root-cert '%s'", cfg.RootCert)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	cli := &GRPCStreamClient{}
	cli.logger = zerolog.New(os.Stdout).With().Str("from", "grpc stream client").Logger()
	cli.cfg = cfg
	if cli.conn, err = grpc.Dial(cfg.Address, opts...); err != nil {
		return nil, errors.Wrapf(err, "failed to create tls-grpc-connection with address %s", cfg.Address)
	}
	cli.client = api.NewGrpcStreamServiceClient(cli.conn)

	return cli, nil
}

// Close 停止运行gRPC流客户端.
func (cli *GRPCStreamClient) Close() {
	if cli.conn != nil {
		cli.conn.Close() // nolint
	}
}

// UploadFile 上传文件.
func (cli *GRPCStreamClient) UploadFile(ctx context.Context, fn string) (*common.Stats, error) {
	var (
		status *api.UploadStatus
		stats  = &common.Stats{}
	)

	fd, err := os.Open(fn)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file '%s'", fn)
	}
	defer fd.Close()

	stream, err := cli.client.Upload(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create upload stream for file %s", fn)
	}
	defer func() {
		stream.CloseSend() // nolint
	}()

	// start to send
	stats.StartedAt = time.Now()

	buffer := make([]byte, cli.cfg.ChunkSize)
WRITE_LOOP:
	for {
		n, err := fd.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break WRITE_LOOP
			}
			return nil, errors.Wrapf(err, "failed unexpectedly while copying from file to buffer")
		}

		if err = stream.Send(&api.FileChunk{
			Content: buffer[:n],
		}); err != nil {
			return nil, errors.Wrapf(err, "failed to send chunk via grpc stream")
		}
	}

	// finish to receive
	stats.FinishedAt = time.Now()

	status, err = stream.CloseAndRecv()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to receive upstream status response")
	}
	if status.Code != api.UploadStatusCode_STATUS_CODE_OK {
		return nil, errors.Errorf("upload failed, msg: %s", status.Message)
	}

	return stats, nil
}
