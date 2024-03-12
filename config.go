package engarde

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var parsedConfig config

type config struct {
	Client clientConfig `yaml:"client"`
	Server serverConfig `yaml:"server"`
}

type clientConfig struct {
	Description        string        `yaml:"description"`
	ListenAddr         string        `yaml:"listenAddr"`
	DstAddr            string        `yaml:"dstAddr"`
	WriteTimeout       time.Duration `yaml:"writeTimeout"`
	IncludedInterfaces []string      `yaml:"includedInterfaces"`
	ExcludedInterfaces []string      `yaml:"excludedInterfaces"`
	DstOverrides       []dstOverride `yaml:"dstOverrides"`
	MTU                int           `yaml:"mtu"`
	UseTeeReader       bool          `yaml:"useTeeReader"`
}

type serverConfig struct {
	Description   string        `yaml:"description"`
	ListenAddr    string        `yaml:"listenAddr"`
	DstAddr       string        `yaml:"dstAddr"`
	WriteTimeout  time.Duration `yaml:"writeTimeout"`
	ClientTimeout int64         `yaml:"clientTimeout"`
	MTU           int           `yaml:"mtu"`
	UseTeeReader  bool          `yaml:"useTeeReader"`
}

type dstOverride struct {
	IfName  string `yaml:"ifName"`
	DstAddr string `yaml:"dstAddr"`
}

func validateAddr(addr string) error {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	_, err = url.Parse(fmt.Sprintf("foo://%s:%s", host, port))

	return err
}

func parseConfig(mode RunMode, configName string) config {
	yamlFile, err := os.ReadFile(configName)
	handleErr(err, fmt.Sprintf("reading config file %s failed", configName))
	err = yaml.Unmarshal(yamlFile, &parsedConfig)
	handleErr(err, "parsing config file failed")

	switch mode {
	case Server:
		if parsedConfig.Server.Description != "" {
			log.Info(parsedConfig.Server.Description)
		}

		handleErr(validateAddr(parsedConfig.Server.ListenAddr), "invalid listenAddr specified")
		handleErr(validateAddr(parsedConfig.Server.DstAddr), "invalid dstAddr specified")

		if parsedConfig.Server.ClientTimeout == 0 {
			parsedConfig.Server.ClientTimeout = 30
		}

		if parsedConfig.Server.WriteTimeout == 0 {
			parsedConfig.Server.WriteTimeout = 10
		}

		if parsedConfig.Server.MTU == 0 {
			parsedConfig.Server.MTU = 1500
		}
	case Client:
		if parsedConfig.Client.Description != "" {
			log.Info(parsedConfig.Client.Description)
		}

		if parsedConfig.Client.MTU == 0 {
			parsedConfig.Client.MTU = 1500
		}

		handleErr(validateAddr(parsedConfig.Client.ListenAddr), "invalid listenAddr specified")
		handleErr(validateAddr(parsedConfig.Client.DstAddr), "invalid dstAddr specified")

		if parsedConfig.Client.WriteTimeout == 0 {
			parsedConfig.Client.WriteTimeout = 10
		}
	}

	return parsedConfig
}
