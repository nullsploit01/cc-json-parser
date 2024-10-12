package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/nullsploit01/cc-json-parser/parser"
	"github.com/spf13/cobra"
)

var jsonParser bool
var runTests bool

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
			var inputFile *os.File

			if len(args) < 1 {
				cmd.PrintErr("Error: A file name is required as an argument.\n")
				cmd.Usage()
				return
			}

			if inputFile == nil {
				file, err := os.Open(args[0])
				if err != nil {
					cmd.PrintErrf("Error reading file: %v\n", err)
					return
				}

				inputFile = file
			}

			currTime := time.Now()
			if jsonParser {
				_, err := RunParser(args[0])
				if err != nil {
					fmt.Println("Error parsing JSON:", err)
					os.Exit(1)
				}

				fmt.Printf("json parsed successfully in %f seconds!\n", time.Since(currTime).Seconds())
			}

		}

		if runTests {
			currTime := time.Now()

			tests_that_should_pass := "./test_data/pass"
			RunTests(false, tests_that_should_pass)

			tests_that_should_fail := "./test_data/fail"
			RunTests(true, tests_that_should_fail)

			fmt.Printf("tests ran successfully in %f seconds!\n", time.Since(currTime).Seconds())
		}
	},
}

func RunParser(filepath string) (interface{}, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	p := parser.NewParser(string(data))
	return p.Parse()
}

func RunTests(shouldFail bool, testFilePath string) {
	err := filepath.Walk(testFilePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			_, err := RunParser(path)
			if err == nil && shouldFail {
				log.Printf("Test failed for file %s: %v", path, err)
				os.Exit(1)
			}

			if !shouldFail && err != nil {
				log.Printf("Test failed for file %s: %v", path, err)
				os.Exit(1)
			}
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&jsonParser, "json-parser", "j", false, "Check if json is valid")
	rootCmd.Flags().BoolVarP(&runTests, "run-tests", "t", false, "run tests against test json files")
}
