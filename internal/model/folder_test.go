package model_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/bootengine/boot/internal/model"
	"github.com/maxatome/go-testdeep/td"
)

func LoadJSONMock(filename string) string {
	content, err := os.ReadFile(fmt.Sprintf("../mocks/%s.json", filename))
	if err != nil {
		panic(fmt.Errorf("failed to load file %s, err: %w", filename, err))
	}

	return string(content)
}

func Test_FolderUnmarshal(t *testing.T) {
	tests := []struct {
		testname string
		input    string
		expected model.Folder
	}{
		{
			testname: "valid - full of simple folder",
			input:    `{"root":["internal", "pkg", "cmd"]}`,
			expected: model.Folder{
				Name: "root",
				Filers: []model.Filer{
					model.Folder{Name: "internal"},
					model.Folder{Name: "pkg"},
					model.Folder{Name: "cmd"},
				},
			},
		},
		{
			testname: "valid - simple folders and simple files",
			input:    `{"root":["internal", "pkg", "main.go"]}`,
			expected: model.Folder{
				Name: "root",
				Filers: model.FolderStruct{
					model.Folder{Name: "internal"},
					model.Folder{Name: "pkg"},
					model.File{Name: "main.go"},
				},
			},
		},
		{
			testname: "valid - only folders - including complex",
			input:    `{"root":["internal", "pkg", {"cmd": ["install", "remove"]}]}`,
			expected: model.Folder{
				Name: "root",
				Filers: model.FolderStruct{
					model.Folder{Name: "internal"},
					model.Folder{Name: "pkg"},
					model.Folder{Name: "cmd", Filers: model.FolderStruct{
						model.Folder{Name: "install"},
						model.Folder{Name: "remove"},
					}},
				},
			},
		},
		{
			testname: "valid - files (simple) and folders (complex)",
			input:    `{"root":["internal", "pkg", {"cmd": ["install", "root.go"]}]}`,
			expected: model.Folder{
				Name: "root",
				Filers: model.FolderStruct{
					model.Folder{Name: "internal"},
					model.Folder{Name: "pkg"},
					model.Folder{Name: "cmd", Filers: model.FolderStruct{
						model.Folder{Name: "install"},
						model.File{Name: "root.go"},
					}},
				},
			},
		},
		{
			testname: "valid - complex files and folders",
			input:    LoadJSONMock("valid_complex_file_folder"),
			expected: model.Folder{
				Name: "root",
				Filers: model.FolderStruct{
					model.Folder{Name: "internal"},
					model.Folder{Name: "pkg"},
					model.Folder{Name: "cmd", Filers: model.FolderStruct{
						model.Folder{Name: "install"},
						model.File{Name: "root.yaml", TempWrapper: &model.TempWrapper{
							TemplateDef: model.TemplateDef{
								Filepath: "./test",
								Engine:   "jinja2",
							},
						}},
					}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			var got model.Folder

			fmt.Printf("%T\n", got)
			err := json.Unmarshal([]byte(tt.input), &got)
			td.CmpNoError(t, err)

			if err == nil {
				td.Cmp(t, got, tt.expected)
			}
		})
	}
}
