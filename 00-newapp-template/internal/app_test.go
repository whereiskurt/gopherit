package internal_test

import (
	"00-newapp-template/internal"
	"00-newapp-template/internal/pkg"
	"testing"
)

func TestGopherCLI(t *testing.T) {
	t.Logf("Testing creation of CLI Application....")

	t.Parallel()

	config := pkg.NewConfig()
	metrics := pkg.NewMetrics(config.Metrics)
	SetupConfig(config)
	internal.NewApp(config, metrics)
}

func SetupConfig(c *pkg.Config) {
	// Test cases are run from the package folder containing the source file.
	c.TemplateFolder = "./../config/template/"
	c.Client.CacheFolder = "./.." + pkg.DefaultClientCacheFolder
	c.VerboseLevel5 = true
	c.VerboseLevel = "5"
}
