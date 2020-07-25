package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/f0m41h4u7/Astaroth/internal/config"
	"github.com/f0m41h4u7/Astaroth/internal/server"
	"github.com/f0m41h4u7/Astaroth/pkg/collector/linux"
)

var (
	n, m    int64
	cfgFile string
)

func init() {
	flag.Int64Var(&n, "n", 5, "sending interval, in seconds")
	flag.Int64Var(&m, "m", 15, "averaging interval, in seconds")
	flag.StringVar(&cfgFile, "config", "./configs/config.json", "path to json configuration file with metrics settings")
}

func main() {
	flag.Parse()
	err := config.InitConfig(cfgFile)
	if err != nil {
		log.Fatal(err)
	}
	_, err = linux.NewCollector(n, m)
	if err != nil {
		log.Fatal(err)
	}

	grpc := server.InitServer()
	defer grpc.Stop()

	sigs := make(chan os.Signal, 1)
	errs := make(chan error, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
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
