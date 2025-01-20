package codegen

import (
	"testing"
)

func TestLabelize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HelloWorld", "hello world"},
		{"labelizeFunction", "labelize function"},
		{"TestCase", "test case"},
		{"simpleTest", "simple test"},
		{"", ""},
		{"Single", "single"},
		{"multipleWordsInOne", "multiple words in one"},
	}

	for _, test := range tests {
		result := labelize(test.input)
		if result != test.expected {
			t.Errorf("labelize(%q) = %q; expected %q", test.input, result, test.expected)
		}
	}
}

func TestUpperFirst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "Hello"},
		{"world", "World"},
		{"", ""},
		{"a", "A"},
		{"A", "A"},
		{"testCase", "TestCase"},
		{"TestCase", "TestCase"},
	}

	for _, test := range tests {
		result := upperFirst(test.input)
		if result != test.expected {
			t.Errorf("upperFirst(%q) = %q; expected %q", test.input, result, test.expected)
		}
	}
}

func TestFilterGenerics(t *testing.T) {
	tests := []struct {
		generics []fieldInfo
		fields   []fieldInfo
		expected []fieldInfo
	}{
		{
			generics: []fieldInfo{
				{Name: "T", Type: "int"},
				{Name: "U", Type: "string"},
			},
			fields: []fieldInfo{
				{Name: "field1", Type: "T"},
				{Name: "field2", Type: "U"},
			},
			expected: []fieldInfo{
				{Name: "T", Type: "int"},
				{Name: "U", Type: "string"},
			},
		},
		{
			generics: []fieldInfo{
				{Name: "T", Type: "int"},
				{Name: "U", Type: "string"},
			},
			fields: []fieldInfo{
				{Name: "field1", Type: "T"},
				{Name: "field2", Type: "V"},
			},
			expected: []fieldInfo{
				{Name: "T", Type: "int"},
			},
		},
		{
			generics: []fieldInfo{
				{Name: "T", Type: "int"},
				{Name: "U", Type: "string"},
			},
			fields: []fieldInfo{
				{Name: "field1", Type: "V"},
				{Name: "field2", Type: "W"},
			},
			expected: []fieldInfo{},
		},
		{
			generics: []fieldInfo{
				{Name: "T", Type: "int"},
			},
			fields: []fieldInfo{
				{Name: "field1", Type: "T"},
				{Name: "field2", Type: "T"},
			},
			expected: []fieldInfo{
				{Name: "T", Type: "int"},
			},
		},
	}

	for _, test := range tests {
		result := filterGenerics(test.generics, test.fields)
		if !equalFieldInfoSlices(result, test.expected) {
			t.Errorf("filterGenerics(%v, %v) = %v; expected %v", test.generics, test.fields, result, test.expected)
		}
	}
}
func equalFieldInfoSlices(a, b []fieldInfo) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
func TestContainsGeneric(t *testing.T) {
	tests := []struct {
		fieldType   string
		genericName string
		expected    bool
	}{
		{"T", "T", false},
		{"List[T]", "T", true},
		{"Map[K,V]", "K", true},
		{"Map[K,V]", "V", true},
		{"Map[K,V]", "T", false},
		{"TypeName", "Name", false},
		{"TypeName", "Type", false},
		{"TypeName", "TypeName", false},
		{"TypeName", "type", false},
		{"TypeName", "name", false},
		{"", "T", false},
		{"T", "", false},
		{"", "", false},
	}

	for _, test := range tests {
		result := containsGeneric(test.fieldType, test.genericName)
		if result != test.expected {
			t.Errorf("containsGeneric(%q, %q) = %v; expected %v", test.fieldType, test.genericName, result, test.expected)
		}
	}
}
func TestIsPartOfLargerIdentifier(t *testing.T) {
	tests := []struct {
		fieldType   string
		genericName string
		expected    bool
	}{
		{"TypeName", "Type", true},
		{"TypeName", "Name", true},
		{"TypeName", "TypeName", true},
		{"TypeName", "type", false},
		{"TypeName", "name", false},
		{"List[T]", "T", false},
		{"Map[K,V]", "K", false},
		{"Map[K,V]", "V", false},
		{"Map[K,V]", "T", false},
		{"", "T", false},
		{"T", "", false},
		{"", "", false},
		{"TypeNameExtra", "Name", true},
		{"TypeNameExtra", "Extra", true},
		{"TypeNameExtra", "TypeName", true},
		{"TypeNameExtra", "TypeNameExtra", true},
	}

	for _, test := range tests {
		result := isPartOfLargerIdentifier(test.fieldType, test.genericName)
		if result != test.expected {
			t.Errorf("isPartOfLargerIdentifier(%q, %q) = %v; expected %v", test.fieldType, test.genericName, result, test.expected)
		}
	}
}
func TestIsIdentifierChar(t *testing.T) {
	tests := []struct {
		input    byte
		expected bool
	}{
		{'a', true},
		{'z', true},
		{'A', true},
		{'Z', true},
		{'0', true},
		{'9', true},
		{'_', true},
		{'-', false},
		{'!', false},
		{' ', false},
		{'$', false},
	}

	for _, test := range tests {
		result := isIdentifierChar(test.input)
		if result != test.expected {
			t.Errorf("isIdentifierChar(%q) = %v; expected %v", test.input, result, test.expected)
		}
	}
}

func TestFieldListOf(t *testing.T) {
	tests := []struct {
		generics []fieldInfo
		expected string
	}{
		{
			generics: []fieldInfo{
				{Name: "T", Type: "int"},
				{Name: "U", Type: "string"},
			},
			expected: "[T int, U string]",
		},
		{
			generics: []fieldInfo{
				{Name: "A", Type: "float64"},
			},
			expected: "[A float64]",
		},
		{
			generics: []fieldInfo{},
			expected: "",
		},
		{
			generics: []fieldInfo{
				{Name: "X", Type: "bool"},
				{Name: "Y", Type: "complex128"},
				{Name: "Z", Type: "byte"},
			},
			expected: "[X bool, Y complex128, Z byte]",
		},
	}

	for _, test := range tests {
		result := fieldListOf(test.generics)
		if result != test.expected {
			t.Errorf("fieldListOf(%v) = %q; expected %q", test.generics, result, test.expected)
		}
	}
}
func TestNameListOf(t *testing.T) {
	tests := []struct {
		generics []fieldInfo
		expected string
	}{
		{
			generics: []fieldInfo{
				{Name: "T", Type: "int"},
				{Name: "U", Type: "string"},
			},
			expected: "[T, U]",
		},
		{
			generics: []fieldInfo{
				{Name: "A", Type: "float64"},
			},
			expected: "[A]",
		},
		{
			generics: []fieldInfo{},
			expected: "",
		},
		{
			generics: []fieldInfo{
				{Name: "X", Type: "bool"},
				{Name: "Y", Type: "complex128"},
				{Name: "Z", Type: "byte"},
			},
			expected: "[X, Y, Z]",
		},
	}

	for _, test := range tests {
		result := nameListOf(test.generics)
		if result != test.expected {
			t.Errorf("nameListOf(%v) = %q; expected %q", test.generics, result, test.expected)
		}
	}
}
func TestTypeListOf(t *testing.T) {
	tests := []struct {
		generics []fieldInfo
		expected string
	}{
		{
			generics: []fieldInfo{
				{Name: "T", Type: "int"},
				{Name: "U", Type: "string"},
			},
			expected: "[int, string]",
		},
		{
			generics: []fieldInfo{
				{Name: "A", Type: "float64"},
			},
			expected: "[float64]",
		},
		{
			generics: []fieldInfo{},
			expected: "",
		},
		{
			generics: []fieldInfo{
				{Name: "X", Type: "bool"},
				{Name: "Y", Type: "complex128"},
				{Name: "Z", Type: "byte"},
			},
			expected: "[bool, complex128, byte]",
		},
	}

	for _, test := range tests {
		result := typeListOf(test.generics)
		if result != test.expected {
			t.Errorf("typeListOf(%v) = %q; expected %q", test.generics, result, test.expected)
		}
	}
}
