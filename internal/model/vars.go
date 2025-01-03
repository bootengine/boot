package model

// Vars is an array of `Var`
// the key is the name of the var inside the workflow
// the value is a `Var` with the type of var and a required flag
type Vars []Var

// ValueType define the type a `Var` can have.
type ValueType string

type Var struct {
	Name     string    `json:"name"`
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
