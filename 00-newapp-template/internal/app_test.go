package internal_test

import (
	"00-newapp-template/internal"
	"00-newapp-template/pkg/config"
	"00-newapp-template/pkg/metrics"
	"os"
	"testing"
)

var m = metrics.NewMetrics()

func TestGopherCLI(t *testing.T) {
	os.Args = []string{"gopherit", "client", "list", "--mode=json"}
	c := config.NewConfig()
	SetupConfig(c)
	app := internal.NewApp(c, m)
	app.InvokeCLI()
}

func TestGopherDefaultIsClient(t *testing.T) {
	os.Args = []string{"gopherit", "list", "--mode=json"}
	c := config.NewConfig()
	SetupConfig(c)
	app := internal.NewApp(c, m)
	app.InvokeCLI()
}

func TestGopherCLIHelp(t *testing.T) {
	os.Args = []string{"gopherit", "--help"}
	c := config.NewConfig()
	SetupConfig(c)
	app := internal.NewApp(c, m)
	app.InvokeCLI()
}

func TestGopherCLIClientHelp(t *testing.T) {
	os.Args = []string{"gopherit", "client", "--help"}
	c := config.NewConfig()
	SetupConfig(c)
	app := internal.NewApp(c, m)
	app.InvokeCLI()
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
