PROJECT                      := github.com/amazingchow/grpc-playground
SRC                          := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
GRPC_CONN_BENCHMARK_TARGETS  := simple-proxy simple-server
GRPC_FILE_TRANSFER_TARGETS   := file-transfer-server file-transfer-client
ALL_TARGETS                  := $(GRPC_CONN_BENCHMARK_TARGETS) $(GRPC_FILE_TRANSFER_TARGETS)

all: build

build: $(ALL_TARGETS)

$(GRPC_CONN_BENCHMARK_TARGETS): $(SRC)
	@go build $(GOMODULEPATH)/$(PROJECT)/grpc-conn-benchmark/$@

$(GRPC_FILE_TRANSFER_TARGETS): $(SRC)
	@go build $(GOMODULEPATH)/$(PROJECT)/grpc-file-transfer-tool/$@

clean:
	rm -f $(ALL_TARGETS)

.PHONY: all build clean
