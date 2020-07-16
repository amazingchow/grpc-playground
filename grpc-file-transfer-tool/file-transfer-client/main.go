package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Usage = "a client tool for file sending via grpc stream"
	app.Version = "v1.0.0"
	app.Commands = []cli.Command{
		{
			Name:   "upload",
			Usage:  "upload a file",
			Action: uploadAction,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "addr",
					Usage: "server's grpc endpoint, e.g. 127.0.0.1:8999",
				},
				&cli.IntFlag{
					Name:  "chunk",
					Usage: "chunk size for every single sending, e.g. 4096",
					Value: 4096,
				},
				&cli.BoolFlag{
					Name:  "compressed",
					Usage: "compress the grpc stream or not",
				},
				&cli.StringFlag{
					Name:  "cert",
					Usage: "root cert file",
				},
				&cli.StringFlag{
					Name:  "file",
					Usage: "file to upload",
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func uploadAction(ctx *cli.Context) (err error) {
	var (
		address    = ctx.String("addr")
		chunkSize  = ctx.Int("chunk")
		compressed = ctx.Bool("compressed")
		rootCert   = ctx.String("cert")
		file       = ctx.String("file")
	)

	cli, err := NewGRPCStreamClient(&GRPCStreamClientCfg{
		Address:    address,
		ChunkSize:  chunkSize,
		Compressed: compressed,
		RootCert:   rootCert,
	})
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	stat, err := cli.UploadFile(context.Background(), file)
	if err != nil {
		panic(err)
	}

	fmt.Printf("used %.2f secs to upload '%s', while chunk size = %d\n", stat.FinishedAt.Sub(stat.StartedAt).Seconds(), file, chunkSize)

	return
}
