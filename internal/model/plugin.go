package model

type PluginType string
type PluginAction string

const (
	ModuleType     PluginType = "module"
	FilerType      PluginType = "filer"
	TempEngineType PluginType = "template_engine"
	VCSType        PluginType = "vcs"

	InitAction              PluginAction = "init"
	InstallLocalDepsAction  PluginAction = "install-local-deps"
	InstallGlobalDepsAction PluginAction = "install-global-deps"
	InstallDevDepsAction    PluginAction = "install-dev-deps"
	CommitAction            PluginAction = "commit"
	PushAction              PluginAction = "push"
	VCSAddAction            PluginAction = "add"
	AddOriginAction         PluginAction = "add-origin"
	CreateFileAction        PluginAction = "create-file"
	CreateFolderAction      PluginAction = "create-folder"
	WriteFileAction         PluginAction = "write-file"
	FormatTemplAction       PluginAction = "apply-template"
)

var (
	Capabilities = map[PluginType][]PluginAction{
		VCSType: {
			InitAction,
			CommitAction,
			PushAction,
			VCSAddAction,
			AddOriginAction,
		},
		ModuleType: {
			InitAction,
			InstallDevDepsAction,
			InstallLocalDepsAction,
			InstallGlobalDepsAction,
		},
		FilerType: {
			CreateFileAction,
			CreateFolderAction,
			WriteFileAction,
		},
		TempEngineType: {
			FormatTemplAction,
		},
	}
)
