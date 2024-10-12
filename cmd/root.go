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
	Short: "CCJP - Another JSON Parser",
	Long: `CCJP (Another JSON Parser) is a command-line tool designed to validate and parse JSON files.
			It supports checking the correctness of JSON files and running predefined tests on JSON datasets. 

			Features:
			- Validate JSON syntax and structure.
			- Run predefined tests to ensure JSON integrity.

			Examples:
			# Validate a single JSON file for correct JSON syntax
			ccjp --json-parser path/to/file.json

			# Run predefined tests (no additional arguments needed)
			ccjp --run-tests

			Usage:
			ccjp [flags]
			Flags:
			-j, --json-parser   Validate a JSON file for syntax and structural correctness.
			-t, --run-tests     Run predefined tests against JSON files located in predefined directories.`,
	Run: func(cmd *cobra.Command, args []string) {
		if jsonParser {
			if len(args) < 1 {
				cmd.PrintErr("Error: A file name is required as an argument for JSON validation.\n")
				cmd.Usage()
				return
			}

			file, err := os.Open(args[0])
			if err != nil {
				cmd.PrintErrf("Error reading file: %v\n", err)
				return
			}
			defer file.Close()

			currTime := time.Now()
			_, err = RunParser(args[0])
			if err != nil {
				fmt.Println("Error parsing JSON:", err)
				os.Exit(1)
			}

			fmt.Printf("JSON parsed successfully in %f seconds!\n", time.Since(currTime).Seconds())
		} else if runTests {
			currTime := time.Now()

			tests_that_should_pass := "./test_data/pass"
			tests_that_should_fail := "./test_data/fail"

			RunTests(false, tests_that_should_pass)
			RunTests(true, tests_that_should_fail)

			fmt.Printf("Tests ran successfully in %f seconds!\n", time.Since(currTime).Seconds())
		} else {
			cmd.Usage()
		}
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&jsonParser, "json-parser", "j", false, "Enable JSON validation mode")
	rootCmd.Flags().BoolVarP(&runTests, "run-tests", "t", false, "Run predefined tests without any additional arguments")
}

func RunParser(filepath string) (interface{}, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	p, err := parser.NewParser(string(data))
	if err != nil {
		return nil, err
	}
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
