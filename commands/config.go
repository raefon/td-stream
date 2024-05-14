package commands

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Set API token/key",
		RunE: func(cmd *cobra.Command, args []string) error {
			apiKey, err := cmd.Flags().GetString("apiKey")
			if err != nil {
				return err
			}

			apiToken, err := cmd.Flags().GetString("apiToken")
			if err != nil {
				return err
			}

			keyPath, err := cmd.Flags().GetString("keyPath")
			if err != nil {
				return err
			}

			viper.Set("apiKey", apiKey)
			viper.Set("apiToken", apiToken)
			viper.Set("keyPath", keyPath)
			viper.WriteConfig()

			return nil
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			log.Print("config updated")
		},
	}
)

func init() {
	configCmd.Flags().String("apiKey", "", "API Key")
	configCmd.MarkFlagRequired("apiKey")
	configCmd.Flags().String("apiToken", "", "API Token")
	configCmd.MarkFlagRequired("apiToken")
	configCmd.Flags().String("serviceUrl", "https://marketplace.tensordock.com/api/v0/client", "Service URL")
	configCmd.Flags().String("keyPath", "~/.ssh/id_rsa", "Path to private key used for authentication")
	rootCmd.AddCommand(configCmd)
}
