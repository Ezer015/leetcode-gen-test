package codegen

import (
	"fmt"
	"go/format"
	"strings"
	"text/template"

	"github.com/Ezer015/leetcode-gen-test/utils"
)

const testTag = "go:generate"

type fieldInfo struct {
	Name string
	Type string
}
type testCaseInfo struct {
	Name string
	Desc string
}

type testFuncData struct {
	FuncName string
	Params   []fieldInfo
	Results  []fieldInfo
	Generics []fieldInfo
}
type testCaseData struct {
	FuncName string
	Cases    []testCaseInfo
}

type testFuncMetadata struct {
	pkgName   string
	testFuncs []testFuncData
}
type testCaseMetadata struct {
	pkgName   string
	testCases []testCaseData
}

const testCaseTemplate = `// Auto-generated test case template for {{.FuncName}}
{{- $paramGenerics := FilterGenerics .Generics .Params}}
{{- $resultGenerics := FilterGenerics .Generics .Results}}
{{- $standardizedFuncName := .FuncName | UpperFirst}}
{{- $testCaseInputTypeName := TestCaseInputTypeNameOf $standardizedFuncName}}
{{- $testCaseOutputTypeName := TestCaseOutputTypeNameOf $standardizedFuncName}}
{{- $testCaseTypeName := TestCaseTypeNameOf $standardizedFuncName}}
var (
/* 	
	_ = {{$testCaseTypeName}}{{TypeListOf .Generics}}{
		input: {{$testCaseInputTypeName}}{{TypeListOf $paramGenerics}}{
			{{- range .Params}}
			{{.Name}}: ...,
			{{- end}}
			},
		output: {{$testCaseOutputTypeName}}{{TypeListOf $resultGenerics}}{
			{{- range .Results}}
			{{.Name}}: ...,
			{{- end}}
		},
	} 
*/
	
)

type {{$testCaseInputTypeName}}{{FieldListOf $paramGenerics}} struct {
    {{- range .Params}}
    {{.Name}} {{.Type}}
    {{- end}}}

type {{$testCaseOutputTypeName}}{{FieldListOf $resultGenerics}} struct {
    {{- range .Results}}
    {{.Name}} {{.Type}}
    {{- end}}}

type {{$testCaseTypeName}}{{FieldListOf .Generics}} struct {
	name   string
	input  {{$testCaseInputTypeName}}{{NameListOf $paramGenerics}}
	output {{$testCaseOutputTypeName}}{{NameListOf $resultGenerics}}
}`

const testTemplate = `// Auto-generated test for {{.FuncName}}
{{- $standardizedFuncName := .FuncName | UpperFirst}}
{{- $results := .Results}}
func Test{{$standardizedFuncName}}(t *testing.T) {
    {{- range $_, $c := .Cases}}
    t.Run("{{.Desc}}", func(t *testing.T) {
        {{- if $.Results}}{{range $i, $r := $.Results}}{{if $i}}, {{end}}{{$r.Name}}{{end}} := {{end}}{{$.FuncName}}({{range $i, $p := $.Params}}{{if $i}}, {{end}}{{$c.Name}}.input.{{$p.Name}}{{end}})
        {{- range $.Results}}
        if {{.Name}} != {{with $c}}{{.Name}}{{end}}.output.{{.Name}} {
            t.Errorf("{{$.FuncName}}() {{.Name}} = %+v, want {{.Name}} = %+v", {{.Name}}, {{with $c}}{{.Name}}{{end}}.output.{{.Name}})
        }
        {{- end}}
    })
    {{- end}}
}`

func GenerateTestCaseTemplates(content []byte) ([]byte, error) {
	tfMetadata, err := extractTestFuncs(content)
	if err != nil {
		return nil, fmt.Errorf("extracting test function: %v", err)
	}

	// Generate test case template
	tmpl, err := template.New("testcase").Funcs(template.FuncMap{
		"UpperFirst":               upperFirst,
		"FilterGenerics":           filterGenerics,
		"TestCaseTypeNameOf":       utils.TestCaseTypeNameOf,
		"TestCaseInputTypeNameOf":  utils.TestCaseInputTypeNameOf,
		"TestCaseOutputTypeNameOf": utils.TestCaseOutputTypeNameOf,
		"FieldListOf":              fieldListOf,
		"NameListOf":               nameListOf,
		"TypeListOf":               typeListOf,
	}).Parse(testCaseTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing test case template: %v", err)
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf(`//go:generate leetcode-gen-test generate --test-case=$GOFILE
package %s

`, tfMetadata.pkgName))
	for _, tf := range tfMetadata.testFuncs {
		var buf strings.Builder
		if err := tmpl.Execute(&buf, tf); err != nil {
			return nil, fmt.Errorf("executing test case template: %v", err)
		}

		formattedCode, err := format.Source([]byte(buf.String()))
		if err != nil {
			return nil, fmt.Errorf("formatting test case template: %v", err)
		}

		result.Write(formattedCode)
	}

	return []byte(result.String()), nil
}

func GenerateTestTemplates(srcContent []byte, testCaseContent []byte) ([]byte, error) {
	tfMetadata, err := extractTestFuncs(srcContent)
	if err != nil {
		return nil, fmt.Errorf("extracting test function: %v", err)
	}
	tcMetadata, err := extractTestCases(testCaseContent)
	if err != nil {
		return nil, fmt.Errorf("extracting test cases: %v", err)
	}
	if tfMetadata.pkgName != tcMetadata.pkgName {
		return nil, fmt.Errorf("package name mismatch: %s != %s", tfMetadata.pkgName, tcMetadata.pkgName)
	}

	// Generate test template
	tmpl, err := template.New("test").Funcs(template.FuncMap{
		"UpperFirst": upperFirst,
	}).Parse(testTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing test template: %v", err)
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf(`
package %s

import "testing"

`, tcMetadata.pkgName))
	for _, tc := range tcMetadata.testCases {
		var params, results []fieldInfo
		for _, tf := range tfMetadata.testFuncs {
			if tf.FuncName == tc.FuncName {
				params = tf.Params
				results = tf.Results
				break
			}
		}
		if len(params) == 0 || len(results) == 0 {
			for _, tf := range tfMetadata.testFuncs {
				if upperFirst(tf.FuncName) == tc.FuncName {
					tc.FuncName = tf.FuncName
					params = tf.Params
					results = tf.Results
					break
				}
			}
		}

		var buf strings.Builder
		if err := tmpl.Execute(&buf, struct {
			FuncName string
			Cases    []testCaseInfo
			Params   []fieldInfo
			Results  []fieldInfo
		}{
			FuncName: tc.FuncName,
			Cases:    tc.Cases,
			Params:   params,
			Results:  results,
		}); err != nil {
			return nil, fmt.Errorf("executing test template: %v", err)
		}

		formattedCode, err := format.Source([]byte(buf.String()))
		if err != nil {
			return nil, fmt.Errorf("formatting test template: %v", err)
		}

		result.Write(formattedCode)
		result.WriteString("\n")
	}

	return []byte(result.String()), nil
}
