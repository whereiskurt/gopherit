package cmd_test

import (
	"00-newapp-template/internal/app/cmd/client"
	pkgclient "00-newapp-template/pkg/client"
	"00-newapp-template/pkg/config"
	"00-newapp-template/pkg/metrics"
	"00-newapp-template/pkg/server"
	"00-newapp-template/pkg/ui"
	"os"
	"testing"
	"time"
)

func TestUnAuthenticatedClient(t *testing.T) {
	serverConfig := config.NewConfig()
	SetupConfig(serverConfig)

	t.Parallel()

	StartServerRunTests(t, ClientTests)

}

func SetupConfig(c *config.Config) {
	c.Server.ListenPort = "10201"
	// Use a different port than the DEFAULT, then we can parallel tests
	c.Client.BaseURL = "http://localhost:10201"
	// Test cases are run from the package folder containing the source file.
	c.TemplateFolder = "./../../../config/template/"
	c.VerboseLevel5 = true
	c.VerboseLevel = "5"

	_ = os.RemoveAll(c.Server.CacheFolder)
	_ = os.RemoveAll(c.Client.CacheFolder)
}

func StartServerRunTests(t *testing.T, f func(*metrics.Metrics, *testing.T)) {
	mm := metrics.NewMetrics()
	c := config.NewConfig()
	SetupConfig(c)
	s := server.NewServer(c, mm)
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
		f(mm, t)
	}
}

func ClientTests(mm *metrics.Metrics, t *testing.T) {
	t.Run("Gopher.List", func(t *testing.T) {
		c := config.NewConfig()
		SetupConfig(c)
		// Show all gopher/things
		c.Client.GopherID = ""
		c.Client.ThingID = ""
		c.Client.SecretKey = ""
		c.Client.AccessKey = ""
		gophers := client.List(pkgclient.NewAdapter(c, mm), ui.NewCLI(c))
		if len(gophers) != 4 {
			t.Errorf("Unexpected count of gophers: %d", len(gophers))
			t.Fail()
		}
	})
	t.Run("Gopher.DELETE.NOAUTH", func(t *testing.T) {
		c := config.NewConfig()
		SetupConfig(c)
		c.Client.GopherID = "1"
		gophers := client.Delete(pkgclient.NewAdapter(c, mm), ui.NewCLI(c))
		if len(gophers) != 1 { // DELETE returns the matching undelete item.
			t.Errorf("Unexpected count of gophers return on UNAUTHORIZED delete: %d - %+v", len(gophers), gophers)
			t.Fail()
		}
	})
	t.Run("Gopher.DELETE.AUTHROIZED", func(t *testing.T) {
		c := config.NewConfig()
		SetupConfig(c)
		c.Client.GopherID = "1"
		c.Client.SecretKey = "notempty"
		c.Client.AccessKey = "notempty"
		gophers := client.Delete(pkgclient.NewAdapter(c, mm), ui.NewCLI(c))
		if len(gophers) != 0 { // DELETE should return empty after successful delete.
			t.Errorf("Unexpected count of gophers after DELETE: %d", len(gophers))
			t.Fail()
		}
	})
}
