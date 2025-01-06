// Package model contains the data representation inside the application
package model

// Workflow is the result of what has been parsed from user's input.
type Workflow struct {
	Config       Config       `json:"config"`
	Vars         Vars         `json:"vars"`
	Steps        []Step       `json:"steps"`
	FolderStruct FolderStruct `json:"folder_struct"`
}

type GeneratingWorkflow struct {
	Config       Config                 `json:"config"`
	Vars         Vars                   `json:"vars"`
	Steps        []Step                 `json:"steps"`
	FolderStruct GeneratingFolderStruct `json:"folder_struct" yaml:"folder_struct"`
}
