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
)

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
	err := config.InitConfig(cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	addr := net.JoinHostPort("127.0.0.1", port)
	grpc := server.InitServer(addr)
	defer grpc.Stop()

	errs := make(chan error, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() { errs <- grpc.Start() }()
	for {
		select {
		case <-sigs:
			signal.Stop(sigs)
			return
		case err = <-errs:
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
