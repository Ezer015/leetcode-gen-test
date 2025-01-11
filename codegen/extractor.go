package codegen

import (
	"go/ast"
	"go/types"
	"strconv"
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
