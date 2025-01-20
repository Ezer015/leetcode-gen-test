package codegen

import (
	"go/ast"
	"go/types"
	"testing"
)

func TestExtractFields(t *testing.T) {
	tests := []struct {
		name   string
		fields *ast.FieldList
		info   *types.Info
		want   []fieldInfo
	}{
		{
			name:   "nil fields",
			fields: nil,
			info:   &types.Info{},
			want:   []fieldInfo{},
		},
		{
			name: "no names",
			fields: &ast.FieldList{
				List: []*ast.Field{
					{Type: &ast.Ident{Name: "int"}},
					{Type: &ast.Ident{Name: "string"}},
				},
			},
			info: &types.Info{
				Types: map[ast.Expr]types.TypeAndValue{
					&ast.Ident{Name: "int"}:    {Type: types.Typ[types.Int]},
					&ast.Ident{Name: "string"}: {Type: types.Typ[types.String]},
				},
			},
			want: []fieldInfo{
				{Name: "field0", Type: "int"},
				{Name: "field1", Type: "string"},
			},
		},
		{
			name: "with names",
			fields: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{{Name: "a"}},
						Type:  &ast.Ident{Name: "int"},
					},
					{
						Names: []*ast.Ident{{Name: "b"}},
						Type:  &ast.Ident{Name: "string"},
					},
				},
			},
			info: &types.Info{
				Types: map[ast.Expr]types.TypeAndValue{
					&ast.Ident{Name: "int"}:    {Type: types.Typ[types.Int]},
					&ast.Ident{Name: "string"}: {Type: types.Typ[types.String]},
				},
			},
			want: []fieldInfo{
				{Name: "a", Type: "int"},
				{Name: "b", Type: "string"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractFields(tt.fields, tt.info); !equalFieldInfoSlices(got, tt.want) {
				t.Errorf("extractFields() = %v, want %v", got, tt.want)
			}
		})
	}
}
