package utils

import (
	"testing"
)

func TestSrcFileNameOf(t *testing.T) {
	tests := []struct {
		testCaseFile string
		expected     string
	}{
		{"example_testcase.go", "example.go"},
		{"test_testcase.go", "test.go"},
		{"sample_testcase.go", "sample.go"},
		{"no_suffix.go", ""},
		{"another_testcase_testcase.go", "another_testcase.go"},
		{"invalid_testcase", ""},
	}

	for _, tt := range tests {
		t.Run(tt.testCaseFile, func(t *testing.T) {
			result := SrcFileNameOf(tt.testCaseFile)
			if result != tt.expected {
				t.Errorf("SrcFileNameOf(%s) = %s; want %s", tt.testCaseFile, result, tt.expected)
			}
		})
	}
}
func TestTestCaseFileNameOf(t *testing.T) {
	tests := []struct {
		sourceFile string
		expected   string
	}{
		{"example.go", "example_testcase.go"},
		{"test.go", "test_testcase.go"},
		{"sample.go", "sample_testcase.go"},
		{"no_suffix.go", "no_suffix_testcase.go"},
		{"another_test.go", "another_test_testcase.go"},
		{"invalid", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.sourceFile, func(t *testing.T) {
			result := TestCaseFileNameOf(tt.sourceFile)
			if result != tt.expected {
				t.Errorf("TestCaseFileNameOf(%s) = %s; want %s", tt.sourceFile, result, tt.expected)
			}
		})
	}
}
func TestTestFileNameOf(t *testing.T) {
	tests := []struct {
		sourceFile string
		expected   string
	}{
		{"example.go", "example_test.go"},
		{"test.go", "test_test.go"},
		{"sample.go", "sample_test.go"},
		{"no_suffix.go", "no_suffix_test.go"},
		{"another_test.go", "another_test_test.go"},
		{"invalid", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.sourceFile, func(t *testing.T) {
			result := TestFileNameOf(tt.sourceFile)
			if result != tt.expected {
				t.Errorf("TestFileNameOf(%s) = %s; want %s", tt.sourceFile, result, tt.expected)
			}
		})
	}
}

func TestTestCaseTypeNameOf(t *testing.T) {
	tests := []struct {
		funcName string
		expected string
	}{
		{"Example", "testExampleCase"},
		{"Test", "testTestCase"},
		{"Sample", "testSampleCase"},
		{"", "testCase"},
		{"AnotherTest", "testAnotherTestCase"},
	}

	for _, tt := range tests {
		t.Run(tt.funcName, func(t *testing.T) {
			result := TestCaseTypeNameOf(tt.funcName)
			if result != tt.expected {
				t.Errorf("TestCaseTypeNameOf(%s) = %s; want %s", tt.funcName, result, tt.expected)
			}
		})
	}
}
func TestTestCaseInputTypeNameOf(t *testing.T) {
	tests := []struct {
		funcName string
		expected string
	}{
		{"Example", "testExampleInput"},
		{"Test", "testTestInput"},
		{"Sample", "testSampleInput"},
		{"", "testInput"},
		{"AnotherTest", "testAnotherTestInput"},
	}

	for _, tt := range tests {
		t.Run(tt.funcName, func(t *testing.T) {
			result := TestCaseInputTypeNameOf(tt.funcName)
			if result != tt.expected {
				t.Errorf("TestCaseInputTypeNameOf(%s) = %s; want %s", tt.funcName, result, tt.expected)
			}
		})
	}
}
func TestTestCaseOutputTypeNameOf(t *testing.T) {
	tests := []struct {
		funcName string
		expected string
	}{
		{"Example", "testExampleOutput"},
		{"Test", "testTestOutput"},
		{"Sample", "testSampleOutput"},
		{"", "testOutput"},
		{"AnotherTest", "testAnotherTestOutput"},
	}

	for _, tt := range tests {
		t.Run(tt.funcName, func(t *testing.T) {
			result := TestCaseOutputTypeNameOf(tt.funcName)
			if result != tt.expected {
				t.Errorf("TestCaseOutputTypeNameOf(%s) = %s; want %s", tt.funcName, result, tt.expected)
			}
		})
	}
}

func TestIsTestCase(t *testing.T) {
	tests := []struct {
		typeName string
		expected bool
	}{
		{"testExampleCase", true},
		{"testTestCase", true},
		{"testSampleCase", true},
		{"ExampleCase", false},
		{"testExample", false},
		{"testExampleInput", false},
		{"testExampleOutput", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.typeName, func(t *testing.T) {
			result := IsTestCase(tt.typeName)
			if result != tt.expected {
				t.Errorf("IsTestCase(%s) = %v; want %v", tt.typeName, result, tt.expected)
			}
		})
	}
}
func TestIsTestCaseInput(t *testing.T) {
	tests := []struct {
		typeName string
		expected bool
	}{
		{"testExampleInput", true},
		{"testTestInput", true},
		{"testSampleInput", true},
		{"ExampleInput", false},
		{"testExample", false},
		{"testExampleCase", false},
		{"testExampleOutput", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.typeName, func(t *testing.T) {
			result := IsTestCaseInput(tt.typeName)
			if result != tt.expected {
				t.Errorf("IsTestCaseInput(%s) = %v; want %v", tt.typeName, result, tt.expected)
			}
		})
	}
}
func TestIsTestCaseOutput(t *testing.T) {
	tests := []struct {
		typeName string
		expected bool
	}{
		{"testExampleOutput", true},
		{"testTestOutput", true},
		{"testSampleOutput", true},
		{"ExampleOutput", false},
		{"testExample", false},
		{"testExampleCase", false},
		{"testExampleInput", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.typeName, func(t *testing.T) {
			result := IsTestCaseOutput(tt.typeName)
			if result != tt.expected {
				t.Errorf("IsTestCase(%s) = %v; want %v", tt.typeName, result, tt.expected)
			}
		})
	}
}

func TestFuncNameOf(t *testing.T) {
	tests := []struct {
		typeName string
		expected string
	}{
		{"testExampleCase", "Example"},
		{"testExampleInput", "Example"},
		{"testExampleOutput", "Example"},
		{"testTestCase", "Test"},
		{"testSampleInput", "Sample"},
		{"testAnotherTestOutput", "AnotherTest"},
		{"ExampleCase", ""},
		{"testExample", ""},
		{"testExampleCaseExtra", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.typeName, func(t *testing.T) {
			result := FuncNameOf(tt.typeName)
			if result != tt.expected {
				t.Errorf("FuncNameOf(%s) = %s; want %s", tt.typeName, result, tt.expected)
			}
		})
	}
}
