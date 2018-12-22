package cmd

import (
	"00-newapp-template/internal/pkg"
	"00-newapp-template/internal/pkg/ui"
	"github.com/spf13/cobra"
)

// Version holds the config and CLI references.
type Version struct {
	Config *pkg.Config
	CLI    *ui.CLI
}

// NewVersion holds a configuration and command line interface reference (for log out, etc.)
func NewVersion(c *pkg.Config) (v Version) {
	v.Config = c
	return
}

// Version just outputs a gopher.
func (v *Version) Version(cmd *cobra.Command, args []string) {
	cli := ui.NewCLI(v.Config)
	cli.DrawGopher()
	return
}
