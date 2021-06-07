package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate AWS temporary credentials when using Virtual MFA",
	Run: func(cmd *cobra.Command, args []string) {
		Generate()
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().Int64VarP(&durationSeconds, "durationSeconds", "d", 43200, "temporary credentials expiration time in seconds")
	generateCmd.Flags().StringVarP(&account, "account", "a", "", "generate temporary credentials for specific account (by default command will generate credentials for all accounts)")
	generateCmd.Flags().StringVarP(&format, "format", "f", "env", "output format, can be \"env\", \"aws\"")
	generateCmd.Flags().StringVarP(&passCode, "passCode", "p", "", "one-time passcode (tool will ask if not specified)")
}

func CheckFormat(format string) {
	switch format {
	case "env":
	case "aws":
	default:
		fmt.Println("Wrong format")
		os.Exit(1)
	}
}

func Generate() {
	c := InitConfig()
	err := c.CheckIfAccountExists()
	cobra.CheckErr(err)

	CheckFormat(format)
	if account == "" {
		for _, a := range c.Accounts {
			s := a.CreateSession()
			a.GetMFADeviceSerial(s)
			if a.MFASerial != "" {
				a.GetTempCredentials(s, &a.MFASerial, durationSeconds, passCode)
				a.Print(format)
			} else {
				fmt.Println("MFA not configured for account:", a.Name)
			}
		}
	} else {
		var e bool
		for _, a := range c.Accounts {
			if account == a.Name {
				e = true
				s := a.CreateSession()
				a.GetMFADeviceSerial(s)
				if a.MFASerial != "" {
					a.GetTempCredentials(s, &a.MFASerial, durationSeconds, passCode)
					a.Print(format)
				} else {
					fmt.Println("MFA not configured for account:", a.Name)
				}
			}
		}
		if !e {
			fmt.Println("Account not found")
		}
	}
}
