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

	c := config.NewConfig()
	m := metrics.NewMetrics()

	c.Server.ListenPort = "20102" // Use a different port than the DEFAULT, then we can parallel tests

	_ = os.RemoveAll(c.Server.CacheFolder)
	_ = os.RemoveAll(c.Client.CacheFolder)

	t.Parallel()

	s := server.NewServer(c, m)

	go func() {
		select {
		case <-time.After(3 * time.Second):
			s.Finished()
		}
	}()

	go func() {
		err := s.ListenAndServe() // BLOCKS
		if err != nil {
			t.Fail()
		}
		s.Finished()
	}()

}
