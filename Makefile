PROJECT                      := github.com/amazingchow/dig-the-grpc
SRC                          := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
CONNECTION_BENCHTEST_TARGETS := server client
FILE_TRANSFER_TARGETS        := file-transfer-server file-transfer-client
ALL_TARGETS                  := $(CONNECTION_BENCHTEST_TARGETS) $(FILE_TRANSFER_TARGETS)

all: build

build: $(ALL_TARGETS)

$(CONNECTION_BENCHTEST_TARGETS): $(SRC)
	@go build $(GOMODULEPATH)/$(PROJECT)/client-connection-benchtest/$@

$(FILE_TRANSFER_TARGETS): $(SRC)
	@go build $(GOMODULEPATH)/$(PROJECT)/grpc-file-transfer-tool/$@

pb-fmt:
	@clang-format -i ./pb/*.proto

clean:
	rm -f $(ALL_TARGETS)

.PHONY: all build clean
