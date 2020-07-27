FROM       golang:1.14.6 AS build
WORKDIR    /astaroth
COPY       . .
RUN        go build -o astaroth ./cmd/astaroth/main.go
ENTRYPOINT ["./astaroth"]
