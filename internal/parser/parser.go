package parser

import (
	_ "embed"
	"fmt"
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

type ParserError struct {
	action, filename string
	err              error
}

func (p ParserError) Error() string {
	return fmt.Sprintf("failed to %s file (%s): %s", p.action, p.filename, p.err)
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
		return nil, ParserError{
			action:   "read",
			err:      err,
			filename: filename,
		}
	}

	if err = p.Check(ctx, *cueValue); err != nil {
		return nil, ParserError{
			action:   "check",
			err:      err,
			filename: filename,
		}

	}

	if err = cueValue.Decode(&workflow); err != nil {
		return nil, ParserError{
			action:   "convert",
			err:      err,
			filename: filename,
		}
	}

	return &workflow, nil
}

func filenameIsURL(filename string) bool {
	reg := regexp.MustCompile("^(http|https)")
	return reg.Match([]byte(filename))
}
