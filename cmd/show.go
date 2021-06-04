package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Checks configuration syntax and outputs its content",
	Run: func(cmd *cobra.Command, args []string) {
		Show()
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func Show() {
	c := InitConfig()
	err := c.CheckIfAccountExists()
	cobra.CheckErr(err)

	out, err := yaml.Marshal(c)
	cobra.CheckErr(err)
	fmt.Println(string(out))
}
