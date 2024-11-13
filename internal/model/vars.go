package model

type Vars map[string]Var

type ValueType string

type Var struct {
	Type     ValueType `json:"type"`
	Required bool      `json:"required"`
}

const (
	String   ValueType = "string"
	License  ValueType = "license"
	Password ValueType = "password"
)
