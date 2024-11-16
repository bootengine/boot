package model

type Step struct {
	Name              string
	Module            string
	Action            ModuleAction
	CurrentWorkingDir string `json:"cwd"`
	Params            []string
}
