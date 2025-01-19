package codegen

import (
	"fmt"
	"strings"
	"unicode"
)

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

func upperFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

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
	return strings.Contains(fieldType, genericName) && !isPartOfLargerIdentifier(fieldType, genericName)
}
func isPartOfLargerIdentifier(fieldType, genericName string) bool {
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
