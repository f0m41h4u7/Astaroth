CMD := astaroth

build:
	go build -o $(CMD) ./cmd/main.go

run:
	make build
	./$(CMD)

test:
	go test -v -count=1 -race -gcflags=-l -timeout=30s ./...

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run ./...

clean:
	rm $(CMD)

gen:
	go generate ./...

.PHONY: test lint clean gen
