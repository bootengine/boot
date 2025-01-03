package model

// ModuleType is a string that contains the type of a module
type ModuleType string

// ModuleAction is a string that represent an action performed by a module
type ModuleAction string

func (p *ModuleType) FromString(s string) {
	*p = ModuleType(s)
}

const (
	CmdType        ModuleType = "cmd"
	FilerType      ModuleType = "filer"
	TempEngineType ModuleType = "template_engine"
	VCSType        ModuleType = "vcs"

	InitAction               ModuleAction = "init"
	InstallLocalDepsAction   ModuleAction = "install-local-deps"
	InstallGlobalDepsAction  ModuleAction = "install-global-deps"
	InstallDevDepsAction     ModuleAction = "install-dev-deps"
	CommitAction             ModuleAction = "commit"
	PushAction               ModuleAction = "push"
	VCSAddAction             ModuleAction = "add"
	AddOriginAction          ModuleAction = "add-origin"
	CreateFileAction         ModuleAction = "create-file"
	CreateFolderAction       ModuleAction = "create-folder"
	CreateFolderStructAction ModuleAction = "create-folder-struct"
	WriteFileAction          ModuleAction = "write-file"
	FormatTemplAction        ModuleAction = "apply-template"
)

var (
	// Capabilities contains the capabilites between [ModuleType] and [ModuleAction].
	// It gives the [ModuleAction] a [ModuleType] is capable of.
	Capabilities = map[ModuleType][]ModuleAction{
		VCSType: {
			InitAction,
			CommitAction,
			PushAction,
			VCSAddAction,
			AddOriginAction,
		},
		CmdType: {
			InitAction,
			InstallDevDepsAction,
			InstallLocalDepsAction,
			InstallGlobalDepsAction,
		},
		FilerType: {
			CreateFileAction,
			CreateFolderAction,
			WriteFileAction,
			CreateFolderStructAction,
		},
		TempEngineType: {
			FormatTemplAction,
		},
	}
)
