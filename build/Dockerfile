FROM       golang:1.14.7-alpine3.12
WORKDIR    /astaroth
COPY       . .
RUN        go build -o astaroth ./cmd/astaroth/main.go
ENTRYPOINT ["./astaroth"]
