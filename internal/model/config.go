package model

// Config define custom behaviour of the runner for the [Workflow] it belongs to.
type Config struct {
	CreateRoot   bool      `json:"create_root" yaml:"create_root"`
	Unrestricted bool      `json:"unrestricted"`
	Includes     []Include `json:"includes,omitempty,omitzero"`
}

// Include represent all the information needed to import another config file
type Include struct {
	From string `json:"from"`
	As   string `json:"as"`
}
