.PHONY: build build-all clean install

build:
	go build -o moka main.go

build-all:
	GOOS=linux GOARCH=amd64 go build -o dist/moka-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -o dist/moka-linux-arm64 main.go
	GOOS=darwin GOARCH=amd64 go build -o dist/moka-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o dist/moka-darwin-arm64 main.go

clean:
	rm -f moka
	rm -rf dist

install: build
	./install.sh
