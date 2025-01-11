package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Ezer015/leetcode-gen-test/codegen"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <source-file>")
		return
	}

	sourceFile := os.Args[1]
	file, err := os.Open(sourceFile)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		return
	}
	defer file.Close()

	// Parse the source file
	testTemplates, err := codegen.GenerateTestTemplates(file)
	if err != nil {
		fmt.Printf("Failed to generate test templates: %v\n", err)
		return
	}

	// Generate test file name
	testFile := fmt.Sprintf("%s_test.go", strings.TrimSuffix(sourceFile, ".go"))

	// Create test file
	f, err := os.Create(testFile)
	if err != nil {
		fmt.Printf("Failed to create test file: %v\n", err)
		return
	}
	defer f.Close()

	// Write test templates to file
	for _, testTemplate := range testTemplates {
		if _, err := f.WriteString(testTemplate); err != nil {
			fmt.Printf("Failed to write test template: %v\n", err)
			return
		}
		f.WriteString("\n")
	}
}
