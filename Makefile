build:
	go build -o /github/go-blockchain/bin/

run: build
	/github/go-blockchain/bin/go-blockchain

test:
	go test -v ./...