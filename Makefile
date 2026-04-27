BINARY=orbx
DIST=dist

.PHONY: build install dev release size lint clean test

build:
	go build -o $(BINARY) .

install:
	go install .

# Local single build (dev)
dev:
	go build -ldflags="-s -w" -o $(BINARY) .

# Full cross-platform release build
release: clean
	@echo "🚀 Building Orbx for all platforms..."

	mkdir -p $(DIST)

	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(DIST)/$(BINARY)-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(DIST)/$(BINARY)-darwin-arm64

	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(DIST)/$(BINARY)-linux-amd64

	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(DIST)/$(BINARY)-windows-amd64.exe

	@echo "✅ Release build complete → $(DIST)/"

size: build
	du -sh $(BINARY)

lint:
	go fmt ./...
	go mod tidy

clean:
	rm -rf $(BINARY) $(DIST)

test:
	go test ./internal/...
