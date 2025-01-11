package codegen

import (
	"os"
	"strings"
	"testing"
)

func TestGenerateTestTemplates(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantErr  bool
		contains []string
	}{
		{
			name: "single return function",
			content: `package example
// @lc
func add(x int, y int) int {
	return x + y
}`,
			wantErr: false,
			contains: []string{
				"package example",
				"func TestAdd",
				"x int",
				"y int",
				"field0 int",
			},
		},
		{
			name: "multiple returns function",
			content: `package example
// @lc
func divide(a int, b int) (int, error) {
	return 0, nil
}`,
			wantErr: false,
			contains: []string{
				"package example",
				"func TestDivide",
				"a int",
				"b int",
				"field0 int",
				"field1 error",
			},
		},
		{
			name: "multiple returns function with named returns",
			content: `package example
// @lc
func swap(a int, b int) (x int, y int) {
	return b, a
}
`,
			wantErr: false,
			contains: []string{
				"package example",
				"func TestSwap",
				"a int",
				"b int",
				"x int",
				"y int",
			},
		},
		{
			name: "no params & return function",
			content: `package example
// @lc
func hello() {
}`,
			wantErr: false,
			contains: []string{
				"package example",
				"func TestHello",
			},
		},
		{
			name: "no lc tag",
			content: `package example
func world() int {

}`,
			wantErr: true,
		},
		{
			name: "no function",
			content: `package example
var x int
`,
			wantErr: true,
		},
		{
			name:    "empty file",
			content: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "test*.go")
			if err != nil {
				t.Errorf("Create Temp File: %v", err)
				return
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.WriteString(tt.content); err != nil {
				t.Fatal(err)
			}
			tmpFile.Seek(0, 0)

			got, err := GenerateTestTemplates(tmpFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateTestTemplates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				for _, want := range tt.contains {
					if !strings.Contains(got[0], want) {
						t.Errorf("GenerateTestTemplates() should contain %q", want)
					}
				}
			}
		})
	}
}
