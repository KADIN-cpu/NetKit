BINARY := netkit
DIST := dist
VERSION := 1.0.0
LDFLAGS := -s -w

.PHONY: build run clean build-all tidy test

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

run: build
	./$(BINARY)

tidy:
	go mod tidy

test:
	go test ./...

clean:
	rm -f $(BINARY)
	rm -rf $(DIST)

build-all: clean
	mkdir -p $(DIST)
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST)/$(BINARY)-windows-amd64.exe .
	GOOS=linux   GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST)/$(BINARY)-linux-amd64 .
	GOOS=linux   GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST)/$(BINARY)-linux-arm64 .
	GOOS=darwin  GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST)/$(BINARY)-darwin-arm64 .
	@echo "Binários gerados em ./$(DIST)"
