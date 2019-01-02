package cmd

import (
	"00-newapp-template/internal/app/cmd/client"
	"00-newapp-template/internal/pkg"
	"00-newapp-template/internal/pkg/adapter"
	"00-newapp-template/internal/pkg/metrics"
	"00-newapp-template/internal/pkg/ui"
	"fmt"
	"github.com/spf13/cobra"
)

// Client is the dispactcher from Cobra to Config
type Client struct {
	Config  *pkg.Config
	Metrics *metrics.Metrics
	Adapter *adapter.Adapter
	CLI     ui.CLI
}

// NewClient dispatches from cobra commands
func NewClient(config *pkg.Config, metrics *metrics.Metrics) (c Client) {
	c.Config = config
	c.Metrics = metrics
	c.CLI = ui.NewCLI(config)
	c.Adapter = adapter.NewAdapter(c.Config, c.Metrics)
	return
}

// Client default action is to show help
func (c *Client) Client(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
	return
}

// List command
func (c *Client) List(cmd *cobra.Command, args []string) {
	client.List(c.Adapter, c.CLI)
	return
}

// Delete command
func (c *Client) Delete(cmd *cobra.Command, args []string) {
	client.Delete(c.Adapter, c.CLI)
	return
}

// Update command
func (c *Client) Update(cmd *cobra.Command, args []string) {
	fmt.Printf("UpdateCommand\n")
	return
}
