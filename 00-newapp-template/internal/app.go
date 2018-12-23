package internal

import (
	"00-newapp-template/internal/app/cmd"
	"00-newapp-template/internal/pkg"
	"00-newapp-template/internal/pkg/ui"
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
	"text/template"
)

// App is created from package main. App handles the configuration and cobra/viper.
type App struct {
	Config  *pkg.Config
	RootCmd *cobra.Command
}

// CommandList entry[0] becomes default when a 'command' is omitted
var CommandList = []string{"client", "server", "version"}

// NewApp constructs the command line and configuration
func NewApp(config * pkg.Config) (a App) {
	a.Config = config
	a.RootCmd = new(cobra.Command)

	// Ensure before any command is run we Unmarshal and Validate the Config values.
	// NOTE: we need to set the PreRun BEFORE making other commands below.
	a.RootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		a.Config.UnmarshalViper()  // copy values from cobra
		err := a.Config.Validate() // and validate.
		if err != nil {
			log.Fatalf("failed to validate: %s", err)
		}
	}

	a.RootCmd.SetUsageTemplate(a.usageTemplate("GopherCLIUsage", nil))
	makeBool("VerboseLevel1", &a.Config.VerboseLevel1, []string{"s", "silent"}, a.RootCmd)
	makeBool("VerboseLevel2", &a.Config.VerboseLevel2, []string{"q", "quiet"}, a.RootCmd)
	makeBool("VerboseLevel3", &a.Config.VerboseLevel3, []string{"v", "info"}, a.RootCmd)
	makeBool("VerboseLevel4", &a.Config.VerboseLevel4, []string{"vv", "debug"}, a.RootCmd)
	makeBool("VerboseLevel5", &a.Config.VerboseLevel5, []string{"vvv", "trace"}, a.RootCmd)
	makeString("VerboseLevel", &a.Config.VerboseLevel, []string{"level"}, a.RootCmd)

	ver := cmd.NewVersion(a.Config)
	vcmd := makeCommand("version", ver.Version, a.RootCmd)
	vcmd.SetUsageTemplate(a.usageTemplate("VersionUsage", nil))
	_ = makeCommand("help", func(command *cobra.Command, i []string) { _ = command.Help() }, vcmd)
	makeBool("Version.ShowServer", &a.Config.Version.ShowServer, []string{"ss", "showserver"}, vcmd)
	makeBool("Version.ShowClient", &a.Config.Version.ShowClient, []string{"sc", "showclient"}, vcmd)

	client := cmd.NewClient(a.Config)
	ccmd := makeCommand("client", client.Client, a.RootCmd)

	makeString("Client.BaseURL", &a.Config.Client.BaseURL, []string{"u", "url"}, ccmd)
	makeString("Client.AccessKey", &a.Config.Client.AccessKey, nil, ccmd)
	makeString("Client.SecretKey", &a.Config.Client.SecretKey, nil, ccmd)
	makeString("Client.CacheKey", &a.Config.Client.CacheKey, nil, ccmd)
	makeString("Client.CacheFolder", &a.Config.Client.CacheFolder, []string{"cf", "cfolder", "cacheFolder"}, ccmd)
	makeBool("Client.CacheResponse", &a.Config.Client.CacheResponse, []string{"c", "cr", "cache", "cacheResponse"}, ccmd)

	makeString("Client.OutputMode", &a.Config.Client.OutputMode, []string{"m", "mode"}, ccmd)
	makeString("Client.GopherID", &a.Config.Client.GopherID, []string{"g", "gid", "gopher", "gopherID"}, ccmd)
	makeString("Client.ThingID", &a.Config.Client.ThingID, []string{"t", "tid", "thing", "thingID"}, ccmd)
	makeString("Client.GopherName", &a.Config.Client.GopherName, []string{"gn", "gname"}, ccmd)
	makeString("Client.GopherDescription", &a.Config.Client.GopherDescription, []string{"gd", "gdesc", "gdescription"}, ccmd)
	makeString("Client.ThingName", &a.Config.Client.ThingName, []string{"tn", "tname"}, ccmd)
	makeString("Client.ThingDescription", &a.Config.Client.ThingDescription, []string{"td", "tdesc", "tdescription"}, ccmd)

	ccmd.SetUsageTemplate(a.usageTemplate("ClientUsage", nil))
	_ = makeCommand("help", func(command *cobra.Command, i []string) { _ = command.Help() }, ccmd)
	_ = makeCommand("list", client.List, ccmd)
	_ = makeCommand("update", client.Update, ccmd)
	_ = makeCommand("delete", client.Delete, ccmd)

	server := cmd.NewServer(a.Config)
	scmd := makeCommand("server", server.Server, a.RootCmd)
	makeString("Server.AccessKey", &a.Config.Server.AccessKey, nil, scmd)
	makeString("Server.SecretKey", &a.Config.Server.SecretKey, nil, scmd)
	makeString("Server.RootFolder", &a.Config.Server.RootFolder, []string{"r", "docroot", "root"}, scmd)
	makeString("Server.ListenPort", &a.Config.Server.ListenPort, []string{"p", "port"}, scmd)

	scmd.SetUsageTemplate(a.usageTemplate("ServerUsage", nil))
	_ = makeCommand("help", func(command *cobra.Command, i []string) { _ = command.Help() }, scmd)
	_ = makeCommand("start", server.Start, scmd)
	_ = makeCommand("stop", server.Stop, scmd)

	return
}

// InvokeCLI passes control over to the root cobra command.
func (a *App) InvokeCLI() {
	// Check if first param & last param if the user wants the Help/Usage
	// NOTE: the "&&" is short-circuited

	arglen := len(os.Args)
	wantsHelp := arglen > 1 && ((os.Args[1] == "--help") || (os.Args[1] == "--h") || (os.Args[1] == "help"))

	if arglen == 1 || wantsHelp {
		usage := a.usageTemplate("GopherCLIUsage", nil)
		_, _ = fmt.Fprintf(os.Stderr, usage)
		return
	}

	// When the CLI is invoked you can do:
	// 		gopherit client list --gopher=1
	//
	// or relying on 'client' as the detault (CommandList[0]):
	// 		gopherit list --gopher=1
	//
	// If the user is relying on the default command then set it.
	setDefaultRootCmd()

	// Call Cobra Execute which will PreRun and select the Command to execute.
	_ = a.RootCmd.Execute()
	return
}

// usageTemplate renders the usage/help/man pages for a cmd
func (a *App) usageTemplate(name string, data interface{}) (usage string) {
	var raw bytes.Buffer
	var err error

	tf := fmt.Sprintf("%scmd/*.tmpl", a.Config.TemplateFolder)

	t := template.New("")
	t, err = t.Funcs(
		template.FuncMap{
			"Gopher": ui.Gopher,
		},
	).ParseGlob(tf)

	if err != nil {
		log.Fatalf("couldn't Template: %v", err)
	}

	err = t.ExecuteTemplate(&raw, name, data)
	if err != nil {
		log.Fatalf("error in Execute template: %v", err)
	}

	usage = raw.String()
	return
}
func setDefaultRootCmd() {
	// Check if the first arg is a root command
	arg := strings.ToLower(os.Args[1])

	// If the first argument isn't one we were expecting, shove CommandList[0] in.
	if !contains(CommandList, arg) {
		// If no root command passed inject the root[0] as default
		rest := os.Args[1:]
		os.Args = []string{os.Args[0], CommandList[0]} // Implant the Default ahead of the rest
		os.Args = append(os.Args, rest...)
	}

	return
}
func contains(a []string, x string) bool {
	for i := range a {
		if x == a[i] {
			return true
		}
	}
	return false
}
func makeCommand(s string, run func(*cobra.Command, []string), parent *cobra.Command) (child *cobra.Command) {
	alias := []string{fmt.Sprintf("%ss", s)} // Add a pluralized alias
	child = &cobra.Command{Use: s, Run: run, PreRun: parent.PreRun, Aliases: alias}
	parent.AddCommand(child)
	return
}
func makeBool(name string, ref *bool, aliases []string, cob *cobra.Command) {
	cob.PersistentFlags().BoolVar(ref, name, *ref, "")
	_ = viper.BindPFlag(name, cob.PersistentFlags().Lookup(name))
	if len(aliases) > 0 && len(aliases[0]) == 1 {
		cob.PersistentFlags().Lookup(name).Shorthand = aliases[0]
	}
	for _, alias := range aliases {
		cob.PersistentFlags().BoolVar(ref, alias, *ref, "")
		cob.PersistentFlags().Lookup(alias).Hidden = true
		viper.RegisterAlias(alias, name)
	}

	return
}
func makeString(name string, ref *string, aliases []string, cob *cobra.Command) {
	cob.PersistentFlags().StringVar(ref, name, *ref, "")
	_ = viper.BindPFlag(name, cob.PersistentFlags().Lookup(name))
	if len(aliases) > 0 && len(aliases[0]) == 1 {
		cob.PersistentFlags().Lookup(name).Shorthand = aliases[0]
	}
	for _, alias := range aliases {
		cob.PersistentFlags().StringVar(ref, alias, *ref, "")
		cob.PersistentFlags().Lookup(alias).Hidden = true
		viper.RegisterAlias(alias, name)
	}

	return
}
