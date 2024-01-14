package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generate documentation for the CLI",
	Long:  `Generate documentation for the CLI`,
	Run: func(cmd *cobra.Command, args []string) {
		output, _ := cmd.Flags().GetString("output")
		if err := doc.GenMarkdownTree(cmd.Parent(), output); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
	docsCmd.PersistentFlags().StringP("output", "o", "./docs", "Output path for docs")
}
