package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/nullsploit01/cc-json-parser/parser"
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
		currTime := time.Now()
		if jsonParser {
			p := parser.NewParser(args[0])
			err := p.Parse()
			if err != nil {
				fmt.Println("Error parsing JSON:", err)
				os.Exit(1)
			}

			fmt.Printf("json parsed successfully in %f seconds!\n", time.Since(currTime).Seconds())
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
