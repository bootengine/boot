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
	InstallLocalDepsAction   ModuleAction = "installLocalDeps"
	InstallGlobalDepsAction  ModuleAction = "installGlobalDeps"
	InstallDevDepsAction     ModuleAction = "installDevDeps"
	CommitAction             ModuleAction = "commit"
	PushAction               ModuleAction = "push"
	VCSAddAction             ModuleAction = "add"
	AddOriginAction          ModuleAction = "addOrigin"
	CreateFileAction         ModuleAction = "createFile"
	CreateFolderAction       ModuleAction = "createFolder"
	CreateFolderStructAction ModuleAction = "createFolderStruct"
	WriteFileAction          ModuleAction = "writeFile"
	FormatTemplAction        ModuleAction = "applyTemplate"
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
