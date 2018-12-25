package server_test

import (
	"00-newapp-template/internal/pkg"
	"00-newapp-template/internal/pkg/server"
	"testing"
	"time"
)

func TestServerStart(t *testing.T) {

	config := pkg.NewConfig()
	config.Server.ListenPort = "20102" // Use a different port than the DEFAULT, then we can parallel tests

	t.Parallel()

	s := server.NewServer(config)

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
