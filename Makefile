PROJECT      := amazingchow/dig-the-grpc
SRC          := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
TARGETS      := server client
ALL_TARGETS  := $(TARGETS)

all: build

build: $(ALL_TARGETS)

$(TARGETS): $(SRC)
	go build $(GOMODULEPATH)/$(PROJECT)/client-connection-benchtest/$@

clean:
	rm -f $(ALL_TARGETS)

.PHONY: all build clean
