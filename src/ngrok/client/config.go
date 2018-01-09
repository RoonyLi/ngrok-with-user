package client

import (
	"fmt"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"net/url"
	"ngrok/log"
	"os"
	"os/user"
	"path"
	"regexp"
	"strings"
)

type Configuration struct {
	HttpProxy          string                          `yaml:"http_proxy,omitempty"`
	ServerAddr         string                          `yaml:"server_addr,omitempty"`
	InspectAddr        string                          `yaml:"inspect_addr,omitempty"`
	TrustHostRootCerts bool                            `yaml:"trust_host_root_certs,omitempty"`
	LogTo              string                          `yaml:"-"`
	Path               string                          `yaml:"-"`
	user               string                          `yaml:"-"`
	password           string                      `yaml:"-"`
	TunnelNames			   []string						   `yaml:"-"`
}


func LoadConfiguration(opts *Options) (config *Configuration, err error) {
	configPath := opts.config
	if configPath == "" {
		configPath = defaultPath()
	}

	log.Info("Reading configuration file %s", configPath)
	configBuf, err := ioutil.ReadFile(configPath)
	if err != nil {
		// failure to read a configuration file is only a fatal error if
		// the user specified one explicitly
		if opts.config != "" {
			err = fmt.Errorf("Failed to read configuration file %s: %v", configPath, err)
			return
		}
	}

	// deserialize/parse the config
	config = new(Configuration)
	if err = yaml.Unmarshal(configBuf, &config); err != nil {
		err = fmt.Errorf("Error parsing configuration file %s: %v", configPath, err)
		return
	}

	// try to parse the old .ngrok format for backwards compatibility
	matched := false
	content := strings.TrimSpace(string(configBuf))
	if matched, err = regexp.MatchString("^[0-9a-zA-Z_\\-!]+$", content); err != nil {
		return
	} else if matched {
		config = &Configuration{AuthToken: content}
	}

	// set configuration defaults
	if config.ServerAddr == "" {
		config.ServerAddr = defaultServerAddr
	}

	// 本地视图
	if config.InspectAddr == "" {
		config.InspectAddr = defaultInspectAddr
	}

	if config.HttpProxy == "" {
		config.HttpProxy = os.Getenv("http_proxy")
	}

	// validate and normalize configuration
	if config.InspectAddr != "disabled" {
		if config.InspectAddr, err = normalizeAddress(config.InspectAddr, "inspect_addr"); err != nil {
			return
		}
	}

	if config.ServerAddr, err = normalizeAddress(config.ServerAddr, "server_addr"); err != nil {
		return
	}

	if config.HttpProxy != "" {
		var proxyUrl *url.URL
		if proxyUrl, err = url.Parse(config.HttpProxy); err != nil {
			return
		} else {
			if proxyUrl.Scheme != "http" && proxyUrl.Scheme != "https" {
				err = fmt.Errorf("Proxy url scheme must be 'http' or 'https', got %v", proxyUrl.Scheme)
				return
			}
		}
	}

	

	// override configuration with command-line options
	config.LogTo = opts.logto
	config.Path = configPath
	if opts.user != "" {
		config.user = opts.user
	}

	if opts.password != "" {
		config.password = opts.password
	}

	switch opts.command {

	case "start-all":	
		return
	// start tunnels
	case "start":
		if len(opts.args) == 0 {
			err = fmt.Errorf("You must specify at least one tunnel to start")
			return
		}
		config.TunnelNames = opts.args;		
	default:
		err = fmt.Errorf("Unknown command: %s", opts.command)
		return
	}

	return
}

func defaultPath() string {
	user, err := user.Current()

	// user.Current() does not work on linux when cross compiling because
	// it requires CGO; use os.Getenv("HOME") hack until we compile natively
	homeDir := os.Getenv("HOME")
	if err != nil {
		log.Warn("Failed to get user's home directory: %s. Using $HOME: %s", err.Error(), homeDir)
	} else {
		homeDir = user.HomeDir
	}

	return path.Join(homeDir, ".ngrok")
}



