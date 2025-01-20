package utils

import (
	"fmt"
	"strings"
)

// SrcFileNameOf takes a test case file name and returns the corresponding source file name.
// It removes the "_testcase.go" suffix from the input file name and appends ".go".
//
// Parameters:
//
//	testCaseFile - the name of the test case file.
//
// Returns:
//
//	The name of the corresponding source file.
func SrcFileNameOf(testCaseFile string) string {
	if !strings.HasSuffix(testCaseFile, "_testcase.go") {
		return ""
	}
	return fmt.Sprintf("%s.go", strings.TrimSuffix(testCaseFile, "_testcase.go"))
}

// TestCaseFileNameOf generates a test case file name based on the provided source file name.
// It removes the ".go" suffix from the source file name and appends "_testcase.go" to it.
//
// Parameters:
//   - sourceFile: The name of the source file as a string.
//
// Returns:
//   - A string representing the generated test case file name.
func TestCaseFileNameOf(sourceFile string) string {
	if !strings.HasSuffix(sourceFile, ".go") {
		return ""
	}
	return fmt.Sprintf("%s_testcase.go", strings.TrimSuffix(sourceFile, ".go"))
}

// TestFileNameOf generates the test file name for a given source file.
// It removes the ".go" suffix from the source file name and appends "_test.go".
//
// Parameters:
//   - sourceFile: The name of the source file.
//
// Returns:
//   - A string representing the test file name.
func TestFileNameOf(sourceFile string) string {
	if !strings.HasSuffix(sourceFile, ".go") {
		return ""
	}
	return fmt.Sprintf("%s_test.go", strings.TrimSuffix(sourceFile, ".go"))
}

const (
	testCasePrefix       = "test"
	testCaseSuffix       = "Case"
	testCaseInputSuffix  = "Input"
	testCaseOutputSuffix = "Output"
)

// TestCaseTypeNameOf generates a test case type name by concatenating a prefix,
// the provided function name, and a suffix.
//
// Parameters:
//   - funcName: The name of the function for which the test case type name is being generated.
//
// Returns:
//
//	A string representing the test case type name.
func TestCaseTypeNameOf(funcName string) string {
	return fmt.Sprintf("%s%s%s", testCasePrefix, funcName, testCaseSuffix)
}

// TestCaseInputTypeNameOf generates a test case input type name by concatenating
// a prefix, the provided function name, and a suffix.
//
// Parameters:
// - funcName: The name of the function for which the test case input type name is generated.
//
// Returns:
// A string representing the test case input type name.
func TestCaseInputTypeNameOf(funcName string) string {
	return fmt.Sprintf("%s%s%s", testCasePrefix, funcName, testCaseInputSuffix)
}

// TestCaseOutputTypeNameOf generates the name for the test case output type
// based on the provided function name. It concatenates a prefix, the function
// name, and a suffix to form the output type name.
//
// Parameters:
//   - funcName: The name of the function for which the test case output type name
//     is being generated.
//
// Returns:
// - A string representing the test case output type name.
func TestCaseOutputTypeNameOf(funcName string) string {
	return fmt.Sprintf("%s%s%s", testCasePrefix, funcName, testCaseOutputSuffix)
}

// IsTestCase checks if the given type name starts with a specific prefix and ends with a specific suffix.
// It returns true if both conditions are met, indicating that the type name represents a test case.
//
// Parameters:
//   - typeName: The name of the type to check.
//
// Returns:
//   - bool: True if the type name is a test case, false otherwise.
func IsTestCase(typeName string) bool {
	return strings.HasPrefix(typeName, testCasePrefix) && strings.HasSuffix(typeName, testCaseSuffix)
}

// IsTestCaseInput checks if the given type name represents a test case input.
// It returns true if the type name starts with the predefined test case prefix
// and ends with the predefined test case input suffix.
//
// Parameters:
//   - typeName: The name of the type to check.
//
// Returns:
//   - bool: True if the type name matches the test case input pattern, false otherwise.
func IsTestCaseInput(typeName string) bool {
	return strings.HasPrefix(typeName, testCasePrefix) && strings.HasSuffix(typeName, testCaseInputSuffix)
}

// IsTestCaseOutput checks if the given typeName represents a test case output.
// It returns true if the typeName starts with the predefined testCasePrefix
// and ends with the predefined testCaseOutputSuffix, otherwise it returns false.
//
// Parameters:
//
//	typeName (string): The name of the type to check.
//
// Returns:
//
//	bool: True if the typeName is a test case output, false otherwise.
func IsTestCaseOutput(typeName string) bool {
	return strings.HasPrefix(typeName, testCasePrefix) && strings.HasSuffix(typeName, testCaseOutputSuffix)
}

// FuncNameOf extracts the function name from a given type name if it follows
// a specific naming convention. The type name must start with a predefined
// prefix and end with one of the predefined suffixes. If the type name does
// not match the expected pattern, an empty string is returned.
//
// Parameters:
//   - typeName: The type name string to process.
//
// Returns:
//   - A string representing the extracted function name, or an empty string
//     if the type name does not follow the expected naming convention.
func FuncNameOf(typeName string) string {
	if !strings.HasPrefix(typeName, testCasePrefix) {
		return ""
	}
	typeName = strings.TrimPrefix(typeName, testCasePrefix)

	for _, suffix := range []string{testCaseSuffix, testCaseInputSuffix, testCaseOutputSuffix} {
		if strings.HasSuffix(typeName, suffix) {
			return strings.TrimSuffix(typeName, suffix)
		}
	}
	return ""
}
