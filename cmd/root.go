package cmd

import (
	"github.com/gothicfann/aws-mfa-auth/pkg/mfa"
	"github.com/spf13/cobra"
)

var (
	cfgFile         string
	durationSeconds int64
	account         string
	format          string
)
var rootCmd = &cobra.Command{
	Use:   "aws-mfa-auth",
	Short: "AWS Virtual MFA Authenticator",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", mfa.FindDefaultCfgFileLocation(), "config file to use")
}

func InitConfig() *mfa.Config {
	if cfgFile == "" {
		cfgFile = mfa.FindDefaultCfgFileLocation()
	}
	cfg := mfa.ReadCfgFile(cfgFile)
	return cfg
}
