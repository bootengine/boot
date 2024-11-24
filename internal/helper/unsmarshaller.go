package helper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/pkg/encoding/json"
	"cuelang.org/go/pkg/encoding/yaml"
)

type SupportedFileType string

const (
	JSON SupportedFileType = "json"
	YAML SupportedFileType = "yaml"
	YML  SupportedFileType = "yml"
)

func CueUnmarshalFile(ctx *cue.Context, filename string) (*cue.Value, error) {
	var (
		exp ast.Expr
	)
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	switch strings.ReplaceAll(filepath.Ext(filename), ".", "") {
	case string(JSON):
		exp, err = json.Unmarshal(fileContent)
		if err != nil {
			return nil, err
		}
	case string(YAML), string(YML):
		exp, err = yaml.Unmarshal(fileContent)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported file type")
	}
	if exp == nil {
		return nil, fmt.Errorf("empty workflow file")
	}
	val := ctx.BuildExpr(exp)
	return &val, nil
}
