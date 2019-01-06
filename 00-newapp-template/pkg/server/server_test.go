package server_test

import (
	"00-newapp-template/pkg/config"
	"00-newapp-template/pkg/metrics"
	"00-newapp-template/pkg/server"
	"os"
	"testing"
	"time"
)

func TestServerStart(t *testing.T) {

	config := config.NewConfig()
	metrics := metrics.NewMetrics()

	config.Server.ListenPort = "20102" // Use a different port than the DEFAULT, then we can parallel tests

	os.RemoveAll(config.Server.CacheFolder)
	os.RemoveAll(config.Client.CacheFolder)

	t.Parallel()

	s := server.NewServer(config, metrics)

	go func() {
		err := s.ListenAndServe() // BLOCKS
		if err != nil {
			t.Fail()
		}
	}()

	select {
	case <-time.After(3 * time.Second):
		s.Finished()
	}

}
