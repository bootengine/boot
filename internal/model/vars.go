package model

// Vars is a map
// the key is the name of the var inside the workflow
// the value is a `Var` with the type of var and a required flag
type Vars map[string]Var

type ValueType string

type Var struct {
	Type     ValueType `json:"type"`
	Required bool      `json:"required"`
}

const (
	String      ValueType = "string"
	License     ValueType = "license"
	Password    ValueType = "password"
	Select      ValueType = "select"
	MultiSelect ValueType = "multi"
)
