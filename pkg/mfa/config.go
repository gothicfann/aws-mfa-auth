package mfa

import (
	"fmt"
	"io/ioutil"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const cfgFileName = ".aws-mfa-auth.yaml"

type Config struct {
	Accounts []Account
}

func FindDefaultCfgFileLocation() string {
	home, err := homedir.Dir()
	cobra.CheckErr(err)
	return fmt.Sprintf("%v/%v", home, cfgFileName)
}

func ReadCfgFile(cfgFile string) *Config {
	bs, err := ioutil.ReadFile(cfgFile)
	cobra.CheckErr(err)
	var c Config
	err = yaml.Unmarshal(bs, &c)
	cobra.CheckErr(err)
	return &c
}

func (c *Config) CheckIfAccountExists() error {
	if len(c.Accounts) == 0 {
		return fmt.Errorf("no accounts described in config file")
	}
	return nil
}
