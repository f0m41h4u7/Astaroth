CMD := astaroth
CLIENT_CMD := client

build:
	go build -o $(CMD) ./cmd/astaroth/main.go

run:
	make build
	./$(CMD)

test:
	go test -v -count=1 -race -gcflags=-l -timeout=30s ./...

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run ./...

clean:
	rm $(CMD) $(CLIENT_CMD)

gen:
	go generate ./...

client:
	go build -o $(CLIENT_CMD) ./cmd/client/main.go

.PHONY: test lint clean gen client
