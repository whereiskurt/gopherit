package internal_test

import (
	"00-newapp-template/internal"
	"00-newapp-template/pkg/config"
	"00-newapp-template/pkg/metrics"
	"00-newapp-template/pkg/server"
	"os"
	"testing"
	"time"
)

var m = metrics.NewMetrics()

func TestGopherCLI(t *testing.T) {
	StartServerRunTests(t)
}

func StartServerRunTests(t *testing.T) {
	c := config.NewConfig()
	SetupConfig(c)
	s := server.NewServer(c, m)
	s.EnableDefaultRouter()
	var err error
	go func() {
		err = s.ListenAndServe() // BLOCKS
	}()

	select {
	// Give the server 1 seconds to fail on startup, before we start client tests.
	case <-time.After(1* time.Second):
		if err != nil {
			t.Logf("Failed: %+v", err)
			t.Fail()
			break
		}
		ClientTests(t)
		s.Finished()
	}
}

func ClientTests(t *testing.T) {

	t.Run("App.ServerHelp", func(t *testing.T) {
		os.Args = []string{"gopherit", "server"}
		c := config.NewConfig()
		SetupConfig(c)
		app := internal.NewApp(c, m)
		app.InvokeCLI()
	})

	t.Run("App.Client.List", func(t *testing.T) {
		c := config.NewConfig()
		SetupConfig(c)
		os.Args = []string{"gopherit", "client", "list", "--mode=json"}
		app := internal.NewApp(c, m)
		app.InvokeCLI()
	})
	t.Run("App.Version.Help", func(t *testing.T) {
		os.Args = []string{"gopherit", "--help"}
		c := config.NewConfig()
		SetupConfig(c)
		app := internal.NewApp(c, m)
		app.InvokeCLI()
	})
	t.Run("App.Version.ClientHelp", func(t *testing.T) {
		os.Args = []string{"gopherit", "client", "--help"}
		c := config.NewConfig()
		SetupConfig(c)
		app := internal.NewApp(c, m)
		app.InvokeCLI()
	})

	t.Run("App.ServerStop", func(t *testing.T) {
		os.Args = []string{"gopherit", "server", "stop"}
		c := config.NewConfig()
		SetupConfig(c)
		app := internal.NewApp(c, m)
		app.InvokeCLI()
	})

}

func SetupConfig(c *config.Config) {
	// Test cases are run from the package folder containing the source file.
	c.TemplateFolder = "./../config/template/"
	c.Client.CacheFolder = "./.." + config.DefaultClientCacheFolder
	c.ConfigFolder = "./../config/"
	c.ConfigFilename = "default.test.gophercli"
	c.VerboseLevel5 = true
	c.VerboseLevel = "5"
	c.ValidateOrFatal()
}
