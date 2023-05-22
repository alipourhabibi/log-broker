package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCMD = &cobra.Command{
	Use:   "test",
	Short: "empjy root command",
	Long:  "empty root command",
	Run:   nil,
}

func Execute() {
	if err := RootCMD.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getConfigfilePath(cmd *cobra.Command) string {
	configFlag := cmd.Flags().Lookup("config")
	if configFlag != nil {
		configFilePath := configFlag.Value.String()
		if configFilePath != "" {
			return configFilePath
		}
	}
	return ""
}
