package model

type TempWorkflow struct {
}

type Workflow struct {
	Config       Config       `json:"config"`
	Vars         Vars         `json:"vars"`
	Steps        []Step       `json:"steps"`
	FolderStruct FolderStruct `json:"folder_struct"`
	From         *string      `json:"from"`
}
