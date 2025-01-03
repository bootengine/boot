package model

// A Step define an action that will be executed in the current [Workflow].
// It has a Name used for logging purpose, it will calls an Action from a installed Module.
// This Action will be run in the CurrentWorkingDir (project_root or "." are default value).
type Step struct {
	Name              string
	Module            string
	Action            ModuleAction
	CurrentWorkingDir string   `json:"cwd,omitempty"`
	Params            []string `json:"params,omitempty"`
}
