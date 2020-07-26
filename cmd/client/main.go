package main

import (
	"context"
	"log"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:1337", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	str, err := api.NewAstarothClient(conn).GetStats(context.Background(), new(empty.Empty))
	if err != nil {
		log.Fatal(err)
	}

	for {
		data, err := str.Recv()
		if err != nil {
			log.Fatal(err)
		}

		if data.String() == "" {
			log.Println("Stream is over")
			return
		}

		log.Printf("Data: %s", data.String())
	}
}
