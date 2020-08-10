CMD := astaroth
CLIENT_CMD := client

build:
	go build -o $(CMD) ./cmd/astaroth/main.go

run:
	make build
	./$(CMD)

docker:
	docker build --tag $(CMD) -f ./build/Dockerfile .

test:
	go test -v -count=1 -race -gcflags=-l -timeout=30s ./...

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run ./...

clean:
	rm $(CMD) $(CLIENT_CMD) logs*

gen:
	go generate ./...

client:
	go build -o $(CLIENT_CMD) ./cmd/client/main.go

bdd:
	set -e ;\
	docker-compose -f ./tests/docker-compose-test.yml up --build -d ;\
	test_status_code=0 ;\
	docker-compose -f ./tests/docker-compose-test.yml run integration_tests go test || test_status_code=$$? ;\
	docker-compose -f ./tests/docker-compose-test.yml down ;\
	exit $$test_status_code ;\

.PHONY: build run docker test lint clean gen client bdd
