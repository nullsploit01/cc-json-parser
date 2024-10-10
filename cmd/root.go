package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var jsonParser bool

var rootCmd = &cobra.Command{
	Use:   "ccjp",
	Short: "Another json parser",
	Long: `A longer description that spans multiple lines and likely contains
				examples and usage of using your application. For example:

				Cobra is a CLI library for Go that empowers applications.
				This application is a tool to generate the needed files
				to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {
		if jsonParser {
			fmt.Println(args[0])
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&jsonParser, "json-parser", "j", false, "Check if json is valid")
}
