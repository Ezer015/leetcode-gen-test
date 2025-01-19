package codegen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"strconv"
	"strings"

	"github.com/Ezer015/leetcode-gen-test/utils"
)

const fieldPrefix = "field"

func extractFields(fields *ast.FieldList, info *types.Info) []fieldInfo {
	params := make([]fieldInfo, 0)
	if fields == nil {
		return params
	}

	for i, field := range fields.List {
		typStr := info.Types[field.Type].Type.String()

		if len(field.Names) > 0 {
			for _, name := range field.Names {
				params = append(params, fieldInfo{
					Name: name.Name,
					Type: typStr,
				})
			}
		} else {
			params = append(params, fieldInfo{
				Name: fieldPrefix + strconv.Itoa(i),
				Type: typStr,
			})
		}
	}
	return params
}

func extractTestFuncs(content []byte) (*testFuncMetadata, error) {
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
	tfMetadata := testFuncMetadata{pkgName: f.Name.Name}
	ast.Inspect(f, func(n ast.Node) bool {
		if decl, ok := n.(*ast.FuncDecl); ok {
			if decl.Doc != nil {
				for _, comment := range decl.Doc.List {
					if strings.Contains(comment.Text, testTag) {
						goto func_extraction
					}
				}
			}
			return true

		func_extraction:
			tfMetadata.testFuncs = append(tfMetadata.testFuncs, testFuncData{
				FuncName: decl.Name.Name,
				Params:   extractFields(decl.Type.Params, info),
				Results:  extractFields(decl.Type.Results, info),
				Generics: extractFields(decl.Type.TypeParams, info),
			})
		}
		return true
	})

	if len(tfMetadata.testFuncs) == 0 {
		return nil, fmt.Errorf("no functions found in leetcode block")
	}
	return &tfMetadata, nil
}

const nameAttrName = "name"

func extractTestCases(content []byte) (*testCaseMetadata, error) {
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

	// Traverse the AST to find variables and their types
	tcMetadata := testCaseMetadata{pkgName: f.Name.Name}
	ast.Inspect(f, func(n ast.Node) bool {
		decl, ok := n.(*ast.GenDecl)
		if !ok || decl.Tok != token.VAR {
			return true
		}

		for _, spec := range decl.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

		var_extraction:
			for i, name := range vs.Names {
				if name == nil {
					continue
				}
				typeStr := ""

				// Get the type of the variable
				if tv, ok := info.Types[vs.Type]; ok {
					if tv.Type != nil {
						typeStr = tv.Type.String()
					}
				}

				// Get the inferred type of the variable
				if typeStr == "" {
					if len(vs.Values) == 0 {
						continue
					}

					if tv, ok := info.Types[vs.Values[0]]; ok {
						if t, ok := tv.Type.(*types.Named); ok {
							if _, isStruct := t.Underlying().(*types.Struct); isStruct {
								typeStr = t.Obj().Name()
							}
						}
					}
				}

				// Check if the variable is a test case
				if typeStr == "" || !utils.IsTestCase(typeStr) {
					continue
				}
				funcStr := utils.FuncNameOf(typeStr)

				tcInfo := testCaseInfo{Name: name.Name}
				if compositeLit, ok := vs.Values[i].(*ast.CompositeLit); ok {
					for _, elt := range compositeLit.Elts {
						if kv, ok := elt.(*ast.KeyValueExpr); ok {
							if key, ok := kv.Key.(*ast.Ident); ok && key.Name == nameAttrName {
								if value, ok := kv.Value.(*ast.BasicLit); ok {
									if value.Kind != token.STRING {
										continue
									}
									tcInfo.Desc = value.Value
									break
								}
							}
						}
					}
				}
				if tcInfo.Desc == "" {
					tcInfo.Desc = labelize(tcInfo.Name)
				}

				for i, tc := range tcMetadata.testCases {
					if tc.FuncName == funcStr {
						tcMetadata.testCases[i].Cases = append(tc.Cases, tcInfo)
						continue var_extraction
					}
				}
				tcMetadata.testCases = append(tcMetadata.testCases, testCaseData{
					FuncName: funcStr,
					Cases:    []testCaseInfo{tcInfo},
				})
			}
		}

		return true
	})

	if len(tcMetadata.testCases) == 0 {
		return nil, fmt.Errorf("no test cases found in leetcode block")
	}
	return &tcMetadata, nil
}
