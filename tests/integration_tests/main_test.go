package main

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

const delay = 10 * time.Second

func TestMain(m *testing.M) {
	log.Printf("wait %s for service availability...", delay)
	time.Sleep(delay)

	status := godog.RunWithOptions("integration", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format:    "pretty",
		Paths:     []string{"features"},
		Randomize: 0,
	})

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func FeatureContext(s *godog.Suite) {
	test := new(statsTest)

	s.BeforeScenario(test.openConnection)
	s.Step(`^I send connect request with N and M parameters$`, test.iSendConnectRequestWithNAndMParameters)
	s.Step(`^I subscribe to server$`, test.iSubscribeToServer)
	s.Step(`^I receive stats every N second$`, test.iReceiveStatsEveryNSecond)
	s.AfterScenario(test.closeConnection)
}
