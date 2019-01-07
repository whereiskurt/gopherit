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
	t.Parallel()

	c := config.NewConfig()
	m := metrics.NewMetrics()

	c.Server.ListenPort = "20102" // Use a different port than the DEFAULT, then we can parallel tests

	// Clean up from last run
	_ = os.RemoveAll(c.Server.CacheFolder)
	_ = os.RemoveAll(c.Client.CacheFolder)

	s := server.NewServer(c, m)

	go func() {
		select {
		case <-time.After(3 * time.Second):
			// After 3 seconds shutdown server properly with success
			s.Finished()
		}
	}()

	go func() {
		// If ListAndServe fails it will fail the test.
		err := s.ListenAndServe() // BLOCKS
		if err != nil {
			t.Fail()
		}
		s.Finished()
	}()

}
