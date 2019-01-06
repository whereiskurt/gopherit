package cmd

import (
	"00-newapp-template/internal/app/cmd/server"
	"00-newapp-template/pkg/config"
	"00-newapp-template/pkg/metrics"
	"00-newapp-template/pkg/ui"
	"github.com/spf13/cobra"
)

// Server holds the config and CLI references.
type Server struct {
	Metrics *metrics.Metrics
	Config  *config.Config
	CLI     ui.CLI
}

// NewServer holds a configuration and command line interface reference (for log out, etc.)
func NewServer(config *config.Config, metrics *metrics.Metrics) (s Server) {
	s.Config = config
	s.CLI = ui.NewCLI(config)
	s.Metrics = metrics
	return
}

// Server with no params will show the help
func (s *Server) Server(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
	return
}

// Start will configure a server and start it.
func (s *Server) Start(cmd *cobra.Command, args []string) {
	server.Start(s.Config, s.Metrics)
	return

}

// Stop will signal the server to stop.
func (s *Server) Stop(cmd *cobra.Command, args []string) {
	server.Stop(s.Config, s.Metrics)
	return
}
