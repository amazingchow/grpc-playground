## Upload files with gRPC

### Server

```shell
./file-transfer-server --port=8999 --cert=cert/cert.pem --key=cert/key.pem
```

### Client

```shell
./file-transfer-client upload --addr=127.0.0.1:8999 --chunk=4096 --compressed=false --cert=cert/cert.pem --file=file.txt
```
