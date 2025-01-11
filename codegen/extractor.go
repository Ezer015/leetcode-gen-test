package codegen

import (
	"go/ast"
	"go/types"
	"strconv"
)

const fieldPrefix = "field"

func extractFields(fields *ast.FieldList, info *types.Info) map[string]string {
	params := make(map[string]string)
	if fields == nil {
		return params
	}

	for i, field := range fields.List {
		typStr := info.Types[field.Type].Type.String()

		if len(field.Names) > 0 {
			for _, name := range field.Names {
				params[name.Name] = typStr
			}
		} else {
			params[fieldPrefix+strconv.Itoa(i)] = typStr
		}
	}
	return params
}
