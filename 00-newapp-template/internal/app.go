package internal

import (
	"00-newapp-template/internal/app/cmd"
	"00-newapp-template/pkg/config"
	"00-newapp-template/pkg/metrics"
	"00-newapp-template/pkg/ui"
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

// App is created from package main. App handles the configuration and cobra/viper.
type App struct {
	Config  *config.Config
	Metrics *metrics.Metrics
	RootCmd *cobra.Command
}

// ApplicationName is referenced for the usage help.
var ApplicationName = "gophercli"

// CommandList entry[0] becomes default when a 'command' is omitted
var CommandList = []string{"client", "server", "version"}

// NewApp constructs the command line and configuration
func NewApp(config *config.Config, mmetrics *metrics.Metrics) (a App) {
	a.Config = config
	a.Metrics = mmetrics
	a.RootCmd = new(cobra.Command)

	// Ensure before any command is run we Unmarshal and Validate the Config values.
	// NOTE: we need to set the PreRun BEFORE making other commands below.
	a.RootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		a.Config.UnmarshalViper()  // copy values from cobra
		a.Config.ValidateOrFatal() // and validate.
	}

	makeBool("VerboseLevel1", &a.Config.VerboseLevel1, []string{"s", "silent"}, a.RootCmd)
	makeBool("VerboseLevel2", &a.Config.VerboseLevel2, []string{"q", "quiet"}, a.RootCmd)
	makeBool("VerboseLevel3", &a.Config.VerboseLevel3, []string{"v", "info"}, a.RootCmd)
	makeBool("VerboseLevel4", &a.Config.VerboseLevel4, []string{"vv", "debug"}, a.RootCmd)
	makeBool("VerboseLevel5", &a.Config.VerboseLevel5, []string{"vvv", "trace"}, a.RootCmd)
	makeString("VerboseLevel", &a.Config.VerboseLevel, []string{"level"}, a.RootCmd)

	ver := cmd.NewVersion(a.Config)
	versionCmd := makeCommand("version", ver.Version, a.RootCmd)
	_ = makeCommand("help", func(command *cobra.Command, i []string) { _ = command.Help() }, versionCmd)
	makeBool("Version.ShowServer", &a.Config.Version.ShowServer, []string{"ss", "showserver"}, versionCmd)
	makeBool("Version.ShowClient", &a.Config.Version.ShowClient, []string{"sc", "showclient"}, versionCmd)

	client := cmd.NewClient(a.Config, a.Metrics)
	clientCmd := makeCommand("client", client.Client, a.RootCmd)
	makeString("client.BaseURL", &a.Config.Client.BaseURL, []string{"u", "url"}, clientCmd)
	makeString("client.AccessKey", &a.Config.Client.AccessKey, nil, clientCmd)
	makeString("client.SecretKey", &a.Config.Client.SecretKey, nil, clientCmd)
	makeString("client.CacheKey", &a.Config.Client.CacheKey, nil, clientCmd)
	makeString("client.CacheFolder", &a.Config.Client.CacheFolder, []string{"cf", "cfolder", "cacheFolder"}, clientCmd)
	makeBool("client.CacheResponse", &a.Config.Client.CacheResponse, []string{"c", "cr", "cache", "cacheResponse"}, clientCmd)
	makeString("client.OutputMode", &a.Config.Client.OutputMode, []string{"m", "mode"}, clientCmd)
	makeString("client.GopherID", &a.Config.Client.GopherID, []string{"g", "gid", "gopher", "gopherID"}, clientCmd)
	makeString("client.ThingID", &a.Config.Client.ThingID, []string{"t", "tid", "thing", "thingID"}, clientCmd)
	makeString("client.GopherName", &a.Config.Client.GopherName, []string{"gn", "gname"}, clientCmd)
	makeString("client.GopherDescription", &a.Config.Client.GopherDescription, []string{"gd", "gdesc", "gdescription"}, clientCmd)
	makeString("client.ThingName", &a.Config.Client.ThingName, []string{"tn", "tname"}, clientCmd)
	makeString("client.ThingDescription", &a.Config.Client.ThingDescription, []string{"td", "tdesc", "tdescription"}, clientCmd)
	makeString("client.ConfigFolder", &a.Config.ConfigFolder, []string{"configFolder"}, clientCmd)
	makeString("client.ConfigFilename", &a.Config.ConfigFilename, []string{"configFile"}, clientCmd)
	makeString("client.TemplateFolder", &a.Config.TemplateFolder, []string{"tfolder", "templateFolder"}, clientCmd)

	_ = makeCommand("help", func(command *cobra.Command, i []string) { _ = command.Help() }, clientCmd)
	_ = makeCommand("list", client.List, clientCmd)
	_ = makeCommand("update", client.Update, clientCmd)
	_ = makeCommand("delete", client.Delete, clientCmd)
	_ = makeCommand("add", client.Add, clientCmd)

	server := cmd.NewServer(a.Config, a.Metrics)
	serverCmd := makeCommand("server", server.Server, a.RootCmd)
	makeString("Server.AccessKey", &a.Config.Server.AccessKey, nil, serverCmd)
	makeString("Server.SecretKey", &a.Config.Server.SecretKey, nil, serverCmd)
	makeString("Server.RootFolder", &a.Config.Server.RootFolder, []string{"r", "docroot", "root"}, serverCmd)
	makeString("Server.ListenPort", &a.Config.Server.ListenPort, []string{"p", "port", "sport"}, serverCmd)
	makeString("Server.MetricsListenPort", &a.Config.Server.MetricsListenPort, []string{"mp", "mport", "metricsport", "metricsPort"}, serverCmd)

	_ = makeCommand("help", func(command *cobra.Command, i []string) { _ = command.Help() }, serverCmd)
	_ = makeCommand("start", server.Start, serverCmd)
	_ = makeCommand("stop", server.Stop, serverCmd)

	// TODO: This block doesn't respect the Config settings. ... need to likely use packr2 :-)
	a.RootCmd.SetUsageTemplate(a.usageTemplate("GopherCLIUsage", nil))
	clientCmd.SetUsageTemplate(a.usageTemplate("ClientUsage", nil))
	serverCmd.SetUsageTemplate(a.usageTemplate("ServerUsage", nil))
	versionCmd.SetUsageTemplate(a.usageTemplate("VersionUsage", nil))

	return
}

// InvokeCLI passes control over to the root cobra command.
func (a *App) InvokeCLI() {
	// Check if first param & last param if the user wants the Help/Usage
	// NOTE: the "&&" is short-circuited

	argLength := len(os.Args)
	wantsHelp := argLength > 1 && ((os.Args[1] == "--help") || (os.Args[1] == "--h") || (os.Args[1] == "help"))

	if argLength == 1 || wantsHelp {
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

	var templateFiles []string
	templateFiles = append(templateFiles, CommandList...)
	templateFiles = append(templateFiles, ApplicationName)

	t := template.New("")
	for _,f := range templateFiles {
		file, err := config.TemplateFolder.Open(fmt.Sprintf("cmd/%s.tmpl", f))
		content, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}

		t, err = t.Funcs(
			template.FuncMap{
				"Gopher": ui.Gopher,
			},
		).Parse(string(content))
	}

	if err != nil {
		log.Printf("couldn't load usage templates from: %v", err)
		return
	}

	err = t.ExecuteTemplate(&raw, name, data)
	if err != nil {
		log.Printf("error execute template for usage: %v", err)
		return
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
