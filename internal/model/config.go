package model

// Config define custom behaviour of the runner for the [Workflow] it belongs to.
type Config struct {
	CreateRoot   bool `json:"create_root" yaml:"create_root"`
	Unrestricted bool `json:"unrestricted"`
}
