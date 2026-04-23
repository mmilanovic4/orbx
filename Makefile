BINARY=orbx

build:
	go build -o $(BINARY) .

install:
	go install .

release:
	go build -ldflags="-s -w" -o $(BINARY) .

size: build
	du -sh $(BINARY)

clean:
	rm -f $(BINARY)
