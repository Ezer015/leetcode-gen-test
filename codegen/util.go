package codegen

import (
	"fmt"
	"strings"
	"unicode"
)

// labelize converts a camelCase or PascalCase string into a space-separated,
// lowercase string. For example, "CamelCase" becomes "camel case".
//
// Parameters:
//
//	s - the input string in camelCase or PascalCase format.
//
// Returns:
//
//	A space-separated, lowercase string.
func labelize(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result = append(result, ' ')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

// upperFirst capitalizes the first letter of the given string.
// If the string is empty, it returns the string as is.
//
// Parameters:
//
//	s - the input string to be modified
//
// Returns:
//
//	A new string with the first letter capitalized, or the original string if it is empty.
func upperFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// filterGenerics filters the given generics based on their presence in the fields.
// It returns a slice of fieldInfo containing only the generics that are found in the fields.
//
// Parameters:
// - generics: A slice of fieldInfo representing the generic types to be filtered.
// - fields: A slice of fieldInfo representing the fields to be checked against.
//
// Returns:
// - A slice of fieldInfo containing the generics that are present in the fields.
func filterGenerics(generics, fields []fieldInfo) []fieldInfo {
	var fieldTypes []string
	for _, field := range fields {
		fieldTypes = append(fieldTypes, field.Type)
	}

	intersection := make([]fieldInfo, 0)
	for _, generic := range generics {
		for _, fieldType := range fieldTypes {
			if containsGeneric(generic.Name, fieldType) {
				intersection = append(intersection, generic)
				break
			}
		}
	}

	return intersection
}
func containsGeneric(fieldType, genericName string) bool {
	if genericName == "" {
		return false
	}
	return strings.Contains(fieldType, genericName) && !isPartOfLargerIdentifier(fieldType, genericName)
}
func isPartOfLargerIdentifier(fieldType, genericName string) bool {
	if genericName == "" {
		return false
	}
	if fieldType == genericName {
		return true
	}
	index := strings.Index(fieldType, genericName)
	if index == -1 {
		return false
	}
	// Check if the character before the generic name is a valid identifier character
	if index > 0 && isIdentifierChar(fieldType[index-1]) {
		return true
	}
	// Check if the character after the generic name is a valid identifier character
	if index+len(genericName) < len(fieldType) && isIdentifierChar(fieldType[index+len(genericName)]) {
		return true
	}
	return false
}
func isIdentifierChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}

// fieldListOf takes a slice of fieldInfo structs and returns a formatted string
// representing the list of fields. Each field is formatted as "Name Type" and
// the entire list is enclosed in square brackets and separated by commas.
// If the input slice is empty, an empty string is returned.
//
// Example:
//
//	input: []fieldInfo{{Name: "id", Type: "int"}, {Name: "name", Type: "string"}}
//	output: "[id int, name string]"
func fieldListOf(generics []fieldInfo) string {
	length := len(generics)
	if length == 0 {
		return ""
	}

	fields := make([]string, len(generics))
	for i, generic := range generics {
		fields[i] = fmt.Sprintf("%s %s", generic.Name, generic.Type)
	}
	return fmt.Sprintf("[%s]", strings.Join(fields, ", "))
}

// nameListOf takes a slice of fieldInfo and returns a formatted string
// representing the names of the fields in the slice. If the slice is empty,
// it returns an empty string. The names are joined by a comma and enclosed
// in square brackets.
//
// Parameters:
//
//	generics []fieldInfo - a slice of fieldInfo structs
//
// Returns:
//
//	string - a formatted string of field names or an empty string if the slice is empty
func nameListOf(generics []fieldInfo) string {
	length := len(generics)
	if length == 0 {
		return ""
	}

	types := make([]string, length)
	for i, field := range generics {
		types[i] = field.Name
	}
	return fmt.Sprintf("[%s]", strings.Join(types, ", "))
}

// typeListOf takes a slice of fieldInfo and returns a string representation
// of the types in the format "[type1, type2, ...]". If the input slice is empty,
// it returns an empty string.
//
// Parameters:
//
//	generics []fieldInfo - A slice of fieldInfo containing type information.
//
// Returns:
//
//	string - A formatted string of the types in the input slice.
func typeListOf(generics []fieldInfo) string {
	length := len(generics)
	if length == 0 {
		return ""
	}

	types := make([]string, length)
	for i, field := range generics {
		types[i] = field.Type
	}
	return fmt.Sprintf("[%s]", strings.Join(types, ", "))
}
