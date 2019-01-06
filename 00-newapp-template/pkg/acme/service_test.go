package acme_test

import (
	"00-newapp-template/pkg/acme"
	"00-newapp-template/pkg/config"
	"00-newapp-template/pkg/metrics"
	"00-newapp-template/pkg/server"
	"os"
	"testing"
	"time"
)

func SetupConfig(c *config.Config) {
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

func StartServerRunTests(t *testing.T, f func(*testing.T, *metrics.Metrics)) {
	// We our own server ports and configs.
	c := config.NewConfig()
	m := metrics.NewMetrics()

	SetupConfig(c)
	s := server.NewServer(c, m)
	s.EnableDefaultRouter()
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
		f(t, m)
	}
}

func ServiceTests(t *testing.T, metrics *metrics.Metrics) {
	c := config.NewConfig()
	SetupConfig(c)

	_ = os.RemoveAll(c.Server.CacheFolder)
	_ = os.RemoveAll(c.Client.CacheFolder)

	t.Run("Service.DELETE.Gopher.NOAUTH", func(t *testing.T) {
		ss := acme.NewService(c.Client.BaseURL, "", "")
		ss.EnableMetrics(metrics)
		gophers := ss.DeleteGopher("1")
		if len(gophers) > 0 {
			t.Logf("Failed: %+v", len(gophers))
			t.Fail()
		}
	})
	t.Run("Service.DELETE.Gopher.AUTHORIZED", func(t *testing.T) {
		ss := acme.NewService(c.Client.BaseURL, "notempty", "notempty")
		ss.EnableMetrics(metrics)
		gophers := ss.DeleteGopher("1")
		if len(gophers) > 3 {
			t.Logf("Failed: %+v", len(gophers))
			t.Fail()
		}
	})
	t.Run("Service.DELETE.Thing.NOAUTH", func(t *testing.T) {
		ss := acme.NewService(c.Client.BaseURL, "", "")
		ss.EnableMetrics(metrics)
		things := ss.DeleteThing("2", "2")
		if things != nil {
			t.Fail()
		}
	})

	t.Run("Service.DELETE.Thing.AUTHORIZED", func(t *testing.T) {
		ss := acme.NewService(c.Client.BaseURL, "notempty", "notempty")
		ss.EnableMetrics(metrics)
		things := ss.DeleteThing("2", "2")
		if len(things) != 2 {
			t.Fail()
		}
	})
}
