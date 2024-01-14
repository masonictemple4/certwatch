package cmd

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "certwatch",
	Short: "certwatch is a simple recurring task to validate ssl certificates via a systemd service",
	Long:  `A simple tool to check SSL certificate validity, including expiration and more either by the CLI or on a recurring basis via a systemd service.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		cfg, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatal(err)
		}

		err = godotenv.Load(cfg)
		if err != nil {
			log.Fatal(err)
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "/etc/env/.certwatch.env", "path to your file containing environment variables")
	// rootCmd.PersistentFlags().StringP("duration", "d", "10", "How frequently should the the task run?")
}
