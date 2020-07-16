#!/bin/bash

run () {
    # chunk size in "1024 (1KB), 2048 (2KB), 4096 (4KB), .. 2097152 (2MB)"
    for shift in $(seq 10 21); do
        for iter in $(seq 1 10); do
            upload "$iter" "$shift"
        done
    done
}

upload () {
    iter=$1
    chunk=$((1 << $2))
    ../../file-transfer-client upload --addr=127.0.0.1:8999 --chunk=$chunk --compressed=false --cert=../cert/cert.pem --file=../fixtures/4G.txt 2>&1 | tee output_$(($iter))_$(($chunk)).log
}

run
