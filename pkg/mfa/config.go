package mfa

import (
	"fmt"
	"io/ioutil"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// cfgFileName is default config name
const cfgFileName = ".aws-mfa-auth.yaml"

// Config describes tool configuration
type Config struct {
	Accounts []Account
}

// FindDefaultCfgFileLocation determines default config location
func FindDefaultCfgFileLocation() string {
	home, err := homedir.Dir()
	cobra.CheckErr(err)
	return fmt.Sprintf("%v/%v", home, cfgFileName)
}

// ReadCfgFile reads config file and yaml unmarshals its content
func ReadCfgFile(cfgFile string) *Config {
	bs, err := ioutil.ReadFile(cfgFile)
	cobra.CheckErr(err)
	var c Config
	err = yaml.Unmarshal(bs, &c)
	cobra.CheckErr(err)
	return &c
}

// CheckIfAccountExists checks if any account exists in config file and returns error value
func (c *Config) CheckIfAccountExists() error {
	if len(c.Accounts) == 0 {
		return fmt.Errorf("no accounts described in config file")
	}
	return nil
}
