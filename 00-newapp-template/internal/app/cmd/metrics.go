package cmd

import (
	"00-newapp-template/internal/pkg"
	"00-newapp-template/internal/pkg/metrics"
	"00-newapp-template/internal/pkg/ui"
	"github.com/spf13/cobra"
)

// Metrics holds the config and CLI references.
type Metrics struct {
	Config   *pkg.Config
	MMetrics *metrics.Metrics
	CLI      ui.CLI
}

// NewMetrics holds a configuration and command line interface reference (for log out, etc.)
func NewMetrics(config *pkg.Config, metrics *metrics.Metrics) (m Metrics) {

	m.Config = config
	m.MMetrics = metrics
	m.CLI = ui.NewCLI(config)
	return
}

// Metrics with no params will show the help
func (m *Metrics) Metrics(cmd *cobra.Command, args []string) {

	return
}
