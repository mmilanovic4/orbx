BINARY=orbx

build:
	go build -o $(BINARY) .

install:
	go install .

release:
	go build -ldflags="-s -w" -o $(BINARY) .

size: build
	du -sh $(BINARY)

lint:
	go fmt ./...
	go mod tidy

clean:
	rm -f $(BINARY)

test:
	go test ./internal/...
