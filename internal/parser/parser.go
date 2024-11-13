package parser

import (
	_ "embed"
	"log"
	"regexp"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/encoding/yaml"
	"github.com/bootengine/boot/internal/model"
)

//go:embed workflow.cue
var schemaFile string

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p Parser) Parse(filename string) (*model.Workflow, error) {
	var (
		workflow model.Workflow
		err      error
		yamlFile *ast.File
		ctx      = cuecontext.New()
		schema   = ctx.CompileString(schemaFile).LookupPath(cue.ParsePath("#Workflow"))
	)

	if filenameIsURL(filename) {
		// clone in temp dir or query the file
	} else {
		yamlFile, err = yaml.Extract(filename, nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	yamlAsCUE := ctx.BuildFile(yamlFile)
	unified := schema.Unify(yamlAsCUE)
	if err := unified.Validate(); err != nil {
		return nil, err
	}

	if err = yamlAsCUE.Decode(&workflow); err != nil {
		return nil, err
	}

	// validate workflow

	return &workflow, nil
}

func filenameIsURL(filename string) bool {
	reg := regexp.MustCompile("^(http|https)")
	return reg.Match([]byte(filename))
}
