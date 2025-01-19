package utils

import (
	"fmt"
	"strings"
)

func SrcFileNameOf(testCaseFile string) string {
	return fmt.Sprintf("%s.go", strings.TrimSuffix(testCaseFile, "_testcase.go"))
}
func TestCaseFileNameOf(sourceFile string) string {
	return fmt.Sprintf("%s_testcase.go", strings.TrimSuffix(sourceFile, ".go"))
}
func TestFileNameOf(sourceFile string) string {
	return fmt.Sprintf("%s_test.go", strings.TrimSuffix(sourceFile, ".go"))
}

const (
	testCasePrefix       = "test"
	testCaseSuffix       = "Case"
	testCaseInputSuffix  = "Input"
	testCaseOutputSuffix = "Output"
)

func TestCaseTypeNameOf(funcName string) string {
	return fmt.Sprintf("%s%s%s", testCasePrefix, funcName, testCaseSuffix)
}
func TestCaseInputTypeNameOf(funcName string) string {
	return fmt.Sprintf("%s%s%s", testCasePrefix, funcName, testCaseInputSuffix)
}
func TestCaseOutputTypeNameOf(funcName string) string {
	return fmt.Sprintf("%s%s%s", testCasePrefix, funcName, testCaseOutputSuffix)
}
func IsTestCase(typeName string) bool {
	return strings.HasPrefix(typeName, testCasePrefix) && strings.HasSuffix(typeName, testCaseSuffix)
}
func IsTestCaseInput(typeName string) bool {
	return strings.HasPrefix(typeName, testCasePrefix) && strings.HasSuffix(typeName, testCaseInputSuffix)
}
func IsTestCaseOutput(typeName string) bool {
	return strings.HasPrefix(typeName, testCasePrefix) && strings.HasSuffix(typeName, testCaseOutputSuffix)
}
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
