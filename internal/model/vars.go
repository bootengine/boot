package model

// Vars is an array of [Var]
type Vars []Var

// ValueType define the type a [Var] can have.
type ValueType string

// Var is a user-defined variable in a [Workflow]. It has a Name, a Type ([ValueType]) and a flag if Required.
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
