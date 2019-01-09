package internal_test

import (
	"00-newapp-template/internal"
	"00-newapp-template/pkg/config"
	"00-newapp-template/pkg/metrics"
	"testing"
)

func TestGopherCLI(t *testing.T) {
	t.Logf("Testing creation of CLI Application....")

	c := config.NewConfig()
	m := metrics.NewMetrics()
	SetupConfig(c)
	internal.NewApp(c, m)
}

func SetupConfig(c *config.Config) {
	// Test cases are run from the package folder containing the source file.
	c.TemplateFolder = "./../config/template/"
	c.Client.CacheFolder = "./.." + config.DefaultClientCacheFolder
	c.VerboseLevel5 = true
	c.VerboseLevel = "5"
	c.ValidateOrFatal()
}
