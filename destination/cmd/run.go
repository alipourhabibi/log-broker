package cmd

import (
	"github.com/alipourhabibi/log-broker/destination/configs"
	"github.com/alipourhabibi/log-broker/destination/services"
	"github.com/spf13/cobra"
)

func init() {
	RootCMD.AddCommand(runCMD)
}

var runCMD = &cobra.Command{
	Use:   "run",
	Short: "Run the application",
	Long:  "Run the application",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		cmd.Flags().String("config", "configs/local_config.yaml", "config file path")
		err := cmd.ParseFlags(args)
		if err != nil {
			return err
		}

		configFilePath := getConfigfilePath(cmd)
		if configFilePath != "" {
			return configs.Confs.Load(configFilePath)
		}
		return nil
	},
	RunE: runCmdE,
}

func runCmdE(cmd *cobra.Command, args []string) error {
	s, err := services.NewServer()
	if err != nil {
		return err
	}
	s.Launch()
	return nil
}
