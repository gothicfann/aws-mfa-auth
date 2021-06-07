package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var awsCmd = &cobra.Command{
	Use:   "aws",
	Short: "Configures AWS shared credentials",
	Run: func(cmd *cobra.Command, args []string) {
		ConfigureAws()
	},
}

func init() {
	configureCmd.AddCommand(awsCmd)
	awsCmd.Flags().Int64VarP(&durationSeconds, "durationSeconds", "d", 43200, "temporary credentials expiration time in seconds")
}

func ConfigureAws() {
	c := InitConfig()
	err := c.CheckIfAccountExists()
	cobra.CheckErr(err)

	home, err := homedir.Dir()
	cobra.CheckErr(err)
	awsDir := home + "/.aws"
	if _, err := os.Stat(awsDir); os.IsNotExist(err) {
		err := os.Mkdir(awsDir, 0700)
		cobra.CheckErr(err)
	}
	awsConfigPath := awsDir + "/config"
	awsCredentialsPath := awsDir + "/credentials"

	os.Remove(awsConfigPath)
	os.Remove(awsCredentialsPath)

	awsConfigFile, err := os.OpenFile(awsConfigPath, os.O_CREATE|os.O_WRONLY, 0600)
	cobra.CheckErr(err)
	defer awsConfigFile.Close()
	awsConfigFile.WriteString("")

	awsCredentialsFile, err := os.OpenFile(awsCredentialsPath, os.O_CREATE|os.O_WRONLY, 0600)
	cobra.CheckErr(err)
	defer awsCredentialsFile.Close()
	awsCredentialsFile.WriteString("")

	for _, a := range c.Accounts {
		awsCredentialsFile.WriteString(a.SprintAws())
		awsConfigFile.WriteString(a.SprintAwsRegion())
		s := a.CreateSession()
		a.GetMFADeviceSerial(s)
		if a.MFASerial != "" {
			a.GetTempCredentials(s, &a.MFASerial, durationSeconds, passCode)
			awsCredentialsFile.WriteString(a.SprintMFAdAws())
			awsConfigFile.WriteString(a.SprintMFAdAwsRegion())
		}
	}
	fmt.Println("All accounts configured successfully")
}
