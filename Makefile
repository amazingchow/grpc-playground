PROJECT                      := github.com/amazingchow/photon-dance-grpc-examples
SRC                          := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
CONNECTION_BENCHMARK_TARGETS := simple-proxy simple-server
FILE_TRANSFER_TARGETS        := file-transfer-server file-transfer-client
ALL_TARGETS                  := $(CONNECTION_BENCHMARK_TARGETS) $(FILE_TRANSFER_TARGETS)

all: build

build: $(ALL_TARGETS)

$(CONNECTION_BENCHMARK_TARGETS): $(SRC)
	@go build $(GOMODULEPATH)/$(PROJECT)/connection-benchmark/$@

$(FILE_TRANSFER_TARGETS): $(SRC)
	@go build $(GOMODULEPATH)/$(PROJECT)/grpc-file-transfer-tool/$@

pb-fmt:
	@clang-format -i ./pb/*.proto

clean:
	rm -f $(ALL_TARGETS)

.PHONY: all build clean
