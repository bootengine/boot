package model

type Step struct {
	Name              string
	Module            string
	Action            PluginAction
	CurrentWorkingDir string `json:"cwd"`
	Params            []string
}
