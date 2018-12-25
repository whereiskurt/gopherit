package cmd

import (
	"00-newapp-template/internal/app/cmd/server"
	"00-newapp-template/internal/pkg"
	"00-newapp-template/internal/pkg/ui"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Server holds the config and CLI references.
type Server struct {
	Config *pkg.Config
	CLI    ui.CLI
}

// NewServer holds a configuration and command line interface reference (for log out, etc.)
func NewServer(config *pkg.Config) (c Server) {
	c.Config = config
	c.CLI = ui.NewCLI(config)
	return
}

// Server with no params will show the help
func (c *Server) Server(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
	return
}

// Start will configure a server and start it.
func (c *Server) Start(cmd *cobra.Command, args []string) {

	l := c.Config.Log.WithFields(log.Fields{
		"docroot": c.Config.Server.RootFolder,
		"cache":   c.Config.Server.CacheFolder,
		"port":    c.Config.Server.ListenPort,
	})

	l.Info("starting server")

	server.Start(c.Config)

	l.Info("server finished")

	return
}

// Stop will signla the server to stop.
func (c *Server) Stop(cmd *cobra.Command, args []string) {
	fmt.Printf("Stop Command\n")
	return
}
