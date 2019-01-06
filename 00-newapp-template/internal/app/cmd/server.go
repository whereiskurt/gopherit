package cmd

import (
	"00-newapp-template/internal/app/cmd/server"
	"00-newapp-template/pkg/config"
	"00-newapp-template/pkg/metrics"
	"00-newapp-template/pkg/ui"
	"fmt"
	"github.com/spf13/cobra"
)

// Server holds the config and CLI references.
type Server struct {
	Metrics *metrics.Metrics
	Config  *config.Config
	CLI     ui.CLI
}

// NewServer holds a configuration and command line interface reference (for log out, etc.)
func NewServer(config *config.Config, metrics *metrics.Metrics) (c Server) {
	c.Config = config
	c.CLI = ui.NewCLI(config)
	c.Metrics = metrics

	return
}

// Server with no params will show the help
func (c *Server) Server(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
	return
}

// Start will configure a server and start it.
func (c *Server) Start(cmd *cobra.Command, args []string) {

	server.Start(c.Config, c.Metrics)

	return
}

// Stop will signal the server to stop.
func (c *Server) Stop(cmd *cobra.Command, args []string) {
	fmt.Printf("Stop Command\n")
	return
}
