#!/bin/bash
# libprotoc 3.12.0
# protoc-gen-go v1.24.0

cd $GOMODULEPATH
for i in $(ls $GOMODULEPATH/github.com/amazingchow/photon-dance-grpc-examples/grpc-file-transfer-tool/pb/*.proto); do
	fn=github.com/amazingchow/photon-dance-grpc-examples/grpc-file-transfer-tool/pb/$(basename "$i")
	echo "compile" $fn
	/usr/local/bin/protoc -I/usr/local/include -I . \
		--go_out=plugins=grpc:. "$fn"
done
