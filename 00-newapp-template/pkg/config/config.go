package config

import (
	"00-newapp-template/pkg/metrics"
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"time"
)

// Config holds all parameters for the application and is structured based on the command hierarchy
type Config struct {
	Context        context.Context
	HomeFolder     string
	HomeFilename   string
	ConfigFolder   string
	ConfigFilename string
	TemplateFolder string
	LogFolder      string
	Log            *log.Logger
	VerboseLevel   string
	VerboseLevel1  bool
	VerboseLevel2  bool
	VerboseLevel3  bool
	VerboseLevel4  bool
	VerboseLevel5  bool

	Client  ClientConfig  // 'gophercli client list'
	Server  ServerConfig  // 'gophercli server start'
	Version VersionConfig // 'gophercli version'
	Metrics MetricsConfig // 'gophercli metrics'
}
type MetricsConfig struct{}

// ClientConfig are all of the params for the Client Command
type ClientConfig struct {
	Config            *Config
	BaseURL           string
	AccessKey         string
	SecretKey         string
	CacheKey          string
	CacheFolder       string
	CacheResponse     bool
	OutputMode        string
	GopherID          string
	GopherName        string
	GopherDescription string
	ThingID           string
	ThingName         string
	ThingDescription  string
	MetricsFolder     string
}

// String allows us to mask the Keys we don't want to reveal
func (c *Config) String() string {
	var safeConfig = new(Config)

	spew.Config.MaxDepth = 2
	spew.Config.DisableMethods = true
	
	// Copy config that was passed
	*safeConfig = *c
	// Overwrite sensitive values with the masked value
	mask := "[**MASKED**]"
	safeConfig.Client.AccessKey = mask
	safeConfig.Client.SecretKey = mask
	safeConfig.Client.CacheKey = mask
	safeConfig.Server.AccessKey = mask
	safeConfig.Server.SecretKey = mask
	safeConfig.Server.CacheKey = mask

	s := spew.Sdump(safeConfig)

	return s
}

// ServerConfig are all of the params for the Client Command
type ServerConfig struct {
	Config            *Config
	ListenPort        string
	AccessKey         string // CSV of allowed AccessKeys
	SecretKey         string // CSV of allowed SecretKeys
	RootFolder        string // Server's document root folder
	CacheKey          string
	CacheFolder       string
	CacheResponse     bool
	MetricsListenPort string
	MetricsFolder     string
}

// VersionConfig are all of the params for the Client Command
type VersionConfig struct {
	Config     *Config
	ShowServer bool
	ShowClient bool
}

// NewConfig returns config that has default values set and is hooked to cobra/viper (if invoked.)
func NewConfig() (config *Config) {
	config = new(Config)
	config.SetToDefaults()
	config.Log = log.New()
	cobra.OnInitialize(func() {
		config.readWithViper()
	})
	config.Context = context.Background()

	// Provide access to config variables - ie. log!
	config.Client.Config = config
	config.Server.Config = config
	config.Version.Config = config

	return
}

func (c *ClientConfig) DumpMetrics() {
	pid := os.Getpid()
	dts := time.Now().Format("20060102150405")
	name := fmt.Sprintf("client.%d.%s.prom", pid, dts)
	file := filepath.Join(".", c.MetricsFolder, name)
	metrics.DumpMetrics(file)
}
func (c *ServerConfig) DumpMetrics() {
	pid := os.Getpid()
	dts := time.Now().Format("20060102150405")
	name := fmt.Sprintf("server.%d.%s.prom", pid, dts)
	file := filepath.Join(".", c.MetricsFolder, name)
	metrics.DumpMetrics(file)
}

func (c *ClientConfig) EnableLogging() {
	filename := c.LogFilename()
	c.Config.SetLogFilename(filename)
}
func (c *ServerConfig) EnableLogging() {
	filename := c.LogFilename()
	c.Config.SetLogFilename(filename)
}

// UnmarshalViper copies all of the cobra/viper config data into our Config struct
// This is the delineation between cobra/viper and using our Config struct.
func (c *Config) UnmarshalViper() {
	// Copy everything from the Viper into our Config
	err := viper.Unmarshal(&c)
	if err != nil {
		log.Fatalf("%s", err)
	}
	return
}

func (c *Config) readWithViper() {
	var err error

	viper.SetConfigType(defaultConfigType)

	viper.AddConfigPath(c.ConfigFolder)
	viper.SetConfigName(c.ConfigFilename)
	err = viper.ReadInConfig()
	if err != nil {
		c.Log.Fatalf("fatal: couldn't read in config: %s", err)
	}

	viper.AddConfigPath(c.HomeFolder)
	viper.SetConfigName(c.HomeFilename)
	err = viper.MergeInConfig()
	if err != nil {
		c.Log.Infof("warning: couldn't load configs from: %s or : %s: %s", c.HomeFolder, c.HomeFilename, err)
	}

	viper.AutomaticEnv()

	return
}
