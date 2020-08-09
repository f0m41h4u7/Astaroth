package main

import (
	"context"
	"errors"
	"log"

	"github.com/cucumber/messages-go/v10"
	"github.com/f0m41h4u7/Astaroth/pkg/api"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

const (
	N = 3
	M = 9
)

var errEmptyStats = errors.New("received empty string")

type statsTest struct {
	conn *grpc.ClientConn
	str  api.Astaroth_GetStatsClient
}

func (t *statsTest) openConnection(*messages.Pickle) {
	var err error
	t.conn, err = grpc.Dial("astaroth:1337", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
}

func (t *statsTest) iSendConnectRequestWithNAndMParameters() error {
	_, err := api.NewAstarothClient(t.conn).Connect(context.Background(), &api.ConnectRequest{SendInterval: N, AverageInterval: M})

	return err
}

func (t *statsTest) iSubscribeToServer() (err error) {
	t.str, err = api.NewAstarothClient(t.conn).GetStats(context.Background(), new(empty.Empty))
	return
}

func (t *statsTest) iReceiveStatsEveryNSecond() error {
	data, err := t.str.Recv()
	if err != nil {
		return err
	}
	if data.String() == "" {
		return errEmptyStats
	}

	return nil
}

func (t *statsTest) closeConnection(*messages.Pickle, error) {
	err := t.conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}
