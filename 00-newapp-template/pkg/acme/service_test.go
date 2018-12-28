package acme_test

import (
	"00-newapp-template/internal/pkg"
	"00-newapp-template/internal/pkg/server"
	"00-newapp-template/pkg/acme"
	"os"
	"testing"
	"time"
)

func SetupConfig(c *pkg.Config) {
	c.Server.ListenPort = "10301"
	// c.Server.CacheFolder = "../../" + pkg.DefaultServerCacheFolder
	c.Client.BaseURL = "http://localhost:10301"
	c.VerboseLevel5 = true
	c.VerboseLevel = "5"
}

func TestService(t *testing.T) {

	t.Parallel()

	StartServerRunTests(t, ServiceTests)
}

func StartServerRunTests(t *testing.T, f func(*testing.T)) {
	// We our own server ports and configs.
	config := pkg.NewConfig()
	SetupConfig(config)
	s := server.NewServer(config)
	s.NewRouter()
	var err error
	go func() {
		err = s.ListenAndServe() // BLOCKS
	}()
	select {
	// Give the server 2 seconds to fail on startup, before we start client tests.
	case <-time.After(2 * time.Second):
		if err != nil {
			t.Logf("Failed: %+v", err)
			t.Fail()
			break
		}
		f(t)
	}
}

func ServiceTests(t *testing.T) {
	config := pkg.NewConfig()
	SetupConfig(config)

	os.RemoveAll(config.Server.CacheFolder)
	os.RemoveAll(config.Client.CacheFolder)

	t.Run("Service.DELETE.Gopher.NOAUTH", func(t *testing.T) {
		ss := acme.NewService(config.Client.BaseURL, "", "")
		gophers := ss.DeleteGopher("1")
		if len(gophers) > 0 {
			t.Logf("Failed: %+v", len(gophers))
			t.Fail()
		}
	})
	t.Run("Service.DELETE.Gopher.AUTHORIZED", func(t *testing.T) {
		ss := acme.NewService(config.Client.BaseURL, "notempty", "notempty")
		gophers := ss.DeleteGopher("1")
		if len(gophers) > 3 {
			t.Logf("Failed: %+v", len(gophers))
			t.Fail()
		}
	})
	t.Run("Service.DELETE.Thing.NOAUTH", func(t *testing.T) {
		ss := acme.NewService(config.Client.BaseURL, "", "")
		things := ss.DeleteThing("2", "2")
		if things != nil {
			t.Fail()
		}
	})

	t.Run("Service.DELETE.Thing.AUTHORIZED", func(t *testing.T) {
		ss := acme.NewService(config.Client.BaseURL, "notempty", "notempty")
		things := ss.DeleteThing("2", "2")
		if len(things) != 2 {
			t.Fail()
		}
	})
}
