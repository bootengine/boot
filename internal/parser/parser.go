package parser

import (
	_ "embed"
	"regexp"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/bootengine/boot/internal/helper"
	"github.com/bootengine/boot/internal/model"
)

//go:embed workflow.cue
var schemaFile string

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p Parser) Check(ctx *cue.Context, value cue.Value) error {
	schema := ctx.CompileString(schemaFile).LookupPath(cue.ParsePath("#Workflow"))
	unified := schema.Unify(value)
	return unified.Validate()
}

func (p Parser) Parse(filename string) (*model.Workflow, error) {
	var (
		workflow model.Workflow
		err      error
		ctx      = cuecontext.New()
	)

	if filenameIsURL(filename) {
		// clone in temp dir or query the file
	}

	cueValue, err := helper.CueUnmarshalFile(ctx, filename)
	if err != nil {
	}

	if err = p.Check(ctx, *cueValue); err != nil {
		return nil, err
	}

	if err = cueValue.Decode(&workflow); err != nil {
		return nil, err
	}

	// validate workflow

	return &workflow, nil
}

func filenameIsURL(filename string) bool {
	reg := regexp.MustCompile("^(http|https)")
	return reg.Match([]byte(filename))
}
