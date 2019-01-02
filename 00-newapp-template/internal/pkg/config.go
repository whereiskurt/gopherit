package pkg

import (
	"context"
	"fmt"
	home "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

// These defaults are needed to configure Viper/Cobra
const defaultConfigType = "yaml"
const defaultConfigFilename = "default.gophercli"
const defaultConfigFolder = "./config/"
const defaultTemplateFolder = "./config/template/"
const defaultHomeFilename = ".gophercli"

// Sensible defaults even with out a configuration file present
const defaultVerboseLevel = "3"
const defaultClientOutputMode = "table"
const defaultServerListenPort = "10101"

const defaultMetricsListenPort = "22222"

// Used by the *_test to the set defaults
// DefaultClientCacheFolder stores default client cache file location
const DefaultClientCacheFolder = "./.cache/client/"
const defaultClientCacheResponse = true

// DefaultServerCacheFolder  stores default server cache file location
const DefaultServerCacheFolder = "./.cache/server/"
const defaultServerCacheResponse = true

// Config holds all parameters for the application and is structured based on the command hierarchy
type Config struct {
	Context        context.Context
	HomeFolder     string
	HomeFilename   string
	ConfigFolder   string
	ConfigFilename string
	TemplateFolder string

	Log *log.Logger

	VerboseLevel  string
	VerboseLevel1 bool
	VerboseLevel2 bool
	VerboseLevel3 bool
	VerboseLevel4 bool
	VerboseLevel5 bool

	Client  ClientConfig
	Server  ServerConfig
	Version VersionConfig
	Metrics MetricsConfig
}
type MetricsConfig struct {
	ListenPort string
}

// ClientConfig are all of the params for the Client Command
type ClientConfig struct {
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
}

// ServerConfig are all of the params for the Client Command
type ServerConfig struct {
	ListenPort        string
	AccessKey         string // CSV of allowed AccessKeys
	SecretKey         string // CSV of allowed SecretKeys
	RootFolder        string // Server's document root folder
	CacheKey          string
	CacheFolder       string
	CacheResponse     bool
	MetricsListenPort string
}

// VersionConfig are all of the params for the Client Command
type VersionConfig struct {
	ShowServer bool
	ShowClient bool
}

// NewConfig returns config that has default values set and is hooked to cobra/viper (if invoked.)
func NewConfig() (config *Config) {
	config = new(Config)
	config.useDefaultValues()
	config.Log = log.New()
	cobra.OnInitialize(func() {
		config.readWithViper()
	})
	config.Context = context.Background()

	return
}

// UnmarshalViper copies all of the cobra/viper config data into our Config struct
// This is the deliniation bertween cobra/ciper and using our Config struct.
func (c *Config) UnmarshalViper() {
	// Copy everything from the Viper into our Config
	err := viper.Unmarshal(&c)
	if err != nil {
		log.Fatalf("%s", err)
	}
	return
}

// Validate the string values inside of the Config after copying from Unmarshal or self-setting.
func (c *Config) Validate() (err error) {
	c.validateVerbosity()
	c.validateOutputMode()

	c.Log.SetFormatter(&log.TextFormatter{})
	c.Log.SetOutput(os.Stdout)

	return
}

func (c *Config) readWithViper() {
	var err error

	viper.SetConfigType(defaultConfigType)

	viper.AddConfigPath(c.ConfigFolder)
	viper.SetConfigName(c.ConfigFilename)
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("fatal: couldn't read in config: %s", err)
	}

	viper.AddConfigPath(c.HomeFolder)
	viper.SetConfigName(c.HomeFilename)
	err = viper.MergeInConfig()
	if err != nil {
		log.Printf("warning: couldn't load configs from: %s or : %s: %s", c.HomeFolder, c.HomeFilename, err)
	}

	viper.AutomaticEnv()

	return
}
func (c *Config) useDefaultValues() {
	c.Client.CacheFolder = DefaultClientCacheFolder
	c.Client.CacheResponse = defaultClientCacheResponse
	c.Server.CacheFolder = DefaultServerCacheFolder
	c.Server.CacheResponse = defaultServerCacheResponse

	c.Server.MetricsListenPort = defaultMetricsListenPort
	c.Client.OutputMode = defaultClientOutputMode
	c.Server.ListenPort = defaultServerListenPort
	c.VerboseLevel = defaultVerboseLevel
	c.ConfigFolder = defaultConfigFolder
	c.ConfigFilename = defaultConfigFilename
	c.TemplateFolder = defaultTemplateFolder

	// Find the User's home folder
	folder, err := home.Dir()
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to detect home directory: %v", err))
	} else {
		c.HomeFolder = folder
	}
	c.HomeFilename = defaultHomeFilename
}
func (c *Config) validateOutputMode() {
	switch strings.ToLower(c.Client.OutputMode) {
	case "csv":
	case "json":
	case "xml":
	case "table":

	default:
		log.Fatalf("invalid OutputMode: '%s'", c.Client.OutputMode)
	}
}
func (c *Config) validateVerbosity() {
	if c.hasVerboseLevel() {
		switch {
		case c.VerboseLevel1:
			c.VerboseLevel = "1"
		case c.VerboseLevel2:
			c.VerboseLevel = "2"
		case c.VerboseLevel3:
			c.VerboseLevel = "3"
		case c.VerboseLevel4:
			c.VerboseLevel = "4"
		case c.VerboseLevel5:
			c.VerboseLevel = "5"
		}
	}

	switch c.VerboseLevel {
	case "5":
		c.VerboseLevel5 = true
		c.Log.SetLevel(log.TraceLevel)
	case "4":
		c.VerboseLevel4 = true
		c.Log.SetLevel(log.DebugLevel)
	case "3":
		c.VerboseLevel3 = true
		c.Log.SetLevel(log.InfoLevel)
	case "2":
		c.VerboseLevel1 = true
		c.Log.SetLevel(log.WarnLevel)
	case "1":
		c.VerboseLevel1 = true
		c.Log.SetLevel(log.ErrorLevel)
	}

	if !c.hasVerboseLevel() {
		log.Fatalf("invalid VerboseLevel: '%s'", c.VerboseLevel)
	}

}
func (c *Config) hasVerboseLevel() bool {
	return c.VerboseLevel1 || c.VerboseLevel2 || c.VerboseLevel3 || c.VerboseLevel4 || c.VerboseLevel5
}
