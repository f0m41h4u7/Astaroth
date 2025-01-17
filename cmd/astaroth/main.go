package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/f0m41h4u7/Astaroth/internal/config"
	"github.com/f0m41h4u7/Astaroth/internal/server"
	"github.com/f0m41h4u7/Astaroth/pkg/collector"
)

const defaultInterval int64 = 1 // default scrape interval, in seconds

var (
	port    string
	cfgFile string
)

func init() {
	flag.StringVar(&port, "port", "1337", "GRPC server port")
	flag.StringVar(&cfgFile, "config", "./configs/config.json", "path to json configuration file with metrics settings")
}

func main() {
	flag.Parse()
	metrics, err := config.ReadConfig(cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	collector := collector.NewCollector(metrics)
	done := make(chan struct{}, 1)

	addr := net.JoinHostPort("0.0.0.0", port)
	grpc := server.InitServer(addr, collector)
	defer grpc.Stop()

	errs := make(chan error, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() { errs <- grpc.Start() }()
	go func() { errs <- collector.Run(defaultInterval, done) }()

	for {
		select {
		case <-sigs:
			signal.Stop(sigs)
			close(done)

			return
		case err = <-errs:
			if err != nil {
				close(done)
				log.Fatal(err)
			}
		}
	}
}
