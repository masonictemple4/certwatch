package cmd

import (
	"fmt"

	"github.com/masoncitemple4/certwatch/internal/validator"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check [host]",
	Short: "The check command will validate the host from the command. (Defaults to port 443.)",
	Long:  `Check the host specified in the command Defaults to port 443.`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		if len(args) < 1 {
			cmd.Help()
		}

		hostname := args[0]

		res := validator.Verify(hostname, port)
		if len(res.Errors) > 0 {
			println("Errors encountered during validation:", len(res.Errors))
			fmt.Printf("Errors: %v\n", res.Errors)
			return
		}

		resStr := res.PreviewString()

		fmt.Printf("\nValidation results for %s:%s\n", hostname, port)
		fmt.Printf("%s\n", resStr)

	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.PersistentFlags().StringP("port", "p", "443", "The port to check the host on.")
}
