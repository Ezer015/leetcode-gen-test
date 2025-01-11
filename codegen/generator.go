package codegen

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"strings"
	"text/template"
)

type testTemplateData struct {
	PackageName string
	FuncName    string
	Params      map[string]string
	Returns     map[string]string
}

const testTemplate = `// Auto-generated test file for {{.FuncName}}
package {{.PackageName}}

import "testing"

type test{{.FuncName | UpperFirst}}Input struct {
    {{- range $param, $type := .Params}}
    	{{$param}} {{$type}}
    {{- end}}}

type test{{.FuncName | UpperFirst}}Result struct {
	{{- range $return, $type := .Returns}}
		{{$return}} {{$type}}
	{{- end}}}

func Test{{.FuncName | UpperFirst}}(t *testing.T) {
    testCases := []struct {
        name   string
        input  test{{.FuncName | UpperFirst}}Input
        result test{{.FuncName | UpperFirst}}Result
    }{
        // TODO: Add test cases
	}
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
			{{- if .Returns}}
				{{- $first := true -}}
				{{- range $return, $_ := .Returns -}}
					{{- if not $first}}, {{end -}}
					{{$return}}
					{{- $first = false -}}
				{{- end}} := 
			{{- end}}{{.FuncName}}({{range $param, $type := .Params}}tc.input.{{$param}},{{end}})
			{{- range $return, $_ := .Returns}}
				if {{$return}} != tc.result.{{$return}} {
					t.Errorf("{{$.FuncName}}() {{$return}} = %+v, want {{$return}} = %+v", {{$return}}, tc.result.{{$return}})
				}
			{{- end}}
        })
    }
}`

func GenerateTestTemplates(file *os.File) ([]string, error) {
	if file == nil {
		return nil, fmt.Errorf("nil file provided")
	}

	// Read file content
	content, err := os.ReadFile(file.Name())
	if err != nil {
		return nil, fmt.Errorf("reading file: %v", err)
	}

	// Parse file content
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parsing file: %v", err)
	}

	// Create a type checker
	conf := types.Config{Importer: nil}
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}

	// Type check the file
	_, err = conf.Check("", fset, []*ast.File{f}, info)
	if err != nil {
		return nil, fmt.Errorf("type checking: %v", err)
	}

	// Traverse the AST to find functions in the leetcode block
	var testTemplates []string
	ast.Inspect(f, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			if funcDecl.Doc != nil {
				for _, comment := range funcDecl.Doc.List {
					if strings.Contains(comment.Text, "@lc") {
						goto template_generation
					}
				}
			}
			return true

		template_generation:
			// Extract function data
			data := testTemplateData{
				PackageName: f.Name.Name,
				FuncName:    funcDecl.Name.Name,
				Params:      extractFields(funcDecl.Type.Params, info),
				Returns:     extractFields(funcDecl.Type.Results, info),
			}

			// Generate test template
			tmpl, err := template.New("test").Funcs(template.FuncMap{
				"UpperFirst": upperFirst,
			}).Parse(testTemplate)
			if err != nil {
				return false
			}

			var buf strings.Builder
			if err := tmpl.Execute(&buf, data); err != nil {
				return false
			}

			formattedCode, err := format.Source([]byte(buf.String()))
			if err != nil {
				return false
			}
			testTemplates = append(testTemplates, string(formattedCode))
		}
		return true
	})

	if len(testTemplates) == 0 {
		return nil, fmt.Errorf("no functions found in leetcode block")
	}

	return testTemplates, nil
}
