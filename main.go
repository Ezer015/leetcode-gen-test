package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/Ezer015/leetcode-gen-test/codegen"
	"github.com/Ezer015/leetcode-gen-test/utils"
)

func main() {
	app := &cli.App{
		Name:  "leetcode-gen-test",
		Usage: "Generate test files for LeetCode solution Go source files",
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "Initialize a test case file for a Go source file",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "force",
						Aliases: []string{"f"},
						Usage:   "Force overwrite existing test case file",
					},
				},
				Action: func(c *cli.Context) error {
					if c.NArg() < 1 {
						return cli.Exit("Usage: leetcode-gen-test init <source_file> [--force]", 1)
					}

					sourceFile := c.Args().Get(0)
					force := c.Bool("force")

					// Generate test case file name
					testCaseFile := utils.TestCaseFileNameOf(sourceFile)
					if !force {
						// Check if test case file already exists
						if _, err := os.Stat(testCaseFile); err == nil {
							fmt.Printf("test case file %s already exists\n", testCaseFile)
							return nil
						}
					}

					// Open source file
					file, err := os.Open(sourceFile)
					if err != nil {
						return cli.Exit(fmt.Errorf("failed to open file: %v", err), 1)
					}
					if file == nil {
						return cli.Exit(fmt.Errorf("nil file provided"), 1)
					}
					defer file.Close()
					// Read file content
					content, err := os.ReadFile(file.Name())
					if err != nil {
						return cli.Exit(fmt.Errorf("failed to read file content: %v", err), 1)
					}

					// Parse the source file
					testCaseTemplates, err := codegen.GenerateTestCaseTemplates(content)
					if err != nil {
						return cli.Exit(fmt.Errorf("failed to generate test case templates: %v", err), 1)
					}
					// Create test case file
					f, err := os.Create(testCaseFile)
					if err != nil {
						return cli.Exit(fmt.Errorf("failed to create test case file: %v", err), 1)
					}
					defer f.Close()
					// Write test case templates to file
					if _, err := f.Write(testCaseTemplates); err != nil {
						return cli.Exit(fmt.Errorf("failed to write test case template: %v", err), 1)
					}
					return nil
				},
			},
			{
				Name:  "generate",
				Usage: "Generate a test file for a Go source file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "test-case",
						Aliases: []string{"t"},
						Usage:   "Specify the test case file",
					},
				},
				Action: func(c *cli.Context) error {
					var (
						sourceFile   string
						testCaseFile string
					)
					testCaseVal := c.String("test-case")
					if testCaseVal == "" {
						if c.NArg() < 1 {
							return cli.Exit("Usage: leetcode-gen-test generate <source_file>", 1)
						}
						sourceFile = c.Args().Get(0)
						testCaseFile = utils.TestCaseFileNameOf(sourceFile)
					} else {
						testCaseFile = testCaseVal
						sourceFile = utils.SrcFileNameOf(testCaseFile)
					}

					// Open source srcFile
					srcFile, err := os.Open(sourceFile)
					if err != nil {
						return cli.Exit(fmt.Errorf("failed to open source file: %v", err), 1)
					}
					if srcFile == nil {
						return cli.Exit(fmt.Errorf("nil file provided"), 1)
					}
					defer srcFile.Close()
					// Open test case file
					tcFile, err := os.Open(testCaseFile)
					if err != nil {
						return cli.Exit(fmt.Errorf("failed to open test case file: %v", err), 1)
					}
					if tcFile == nil {
						return cli.Exit(fmt.Errorf("nil file provided"), 1)
					}
					defer tcFile.Close()

					// Read file srcContent
					srcContent, err := os.ReadFile(srcFile.Name())
					if err != nil {
						return cli.Exit(fmt.Errorf("failed to read file content: %v", err), 1)
					}
					// Read file testCaseContent
					testCaseContent, err := os.ReadFile(tcFile.Name())
					if err != nil {
						return cli.Exit(fmt.Errorf("failed to read file content: %v", err), 1)
					}

					// Parse the file
					testTemplates, err := codegen.GenerateTestTemplates(srcContent, testCaseContent)
					if err != nil {
						return cli.Exit(fmt.Errorf("failed to generate test templates: %v", err), 1)
					}
					// Generate test file name
					testFile := utils.TestFileNameOf(sourceFile)
					// Create test file
					f, err := os.Create(testFile)
					if err != nil {
						return cli.Exit(fmt.Errorf("failed to create test file: %v", err), 1)
					}
					defer f.Close()
					// Write test templates to file
					if _, err := f.Write(testTemplates); err != nil {
						return cli.Exit(fmt.Errorf("failed to write test template: %v", err), 1)
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
