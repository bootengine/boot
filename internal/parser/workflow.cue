#Config : {
	create_root?: bool | true
	unrestricted?: bool | false
}

#Var : {
	name: string
	type: "string" | "license" | "password"
	required: bool
}

#Vars: [...#Var]


#StepAction: "init" | "installLocalDeps" | "installGlobalDeps" | "installDevDeps" | "commit"| "push"| "add"| "addOrigin"| "createFile"| "createFolder"| "writeFile"| "applyTemplate" | "createFolderStruct"

#Step : {
	name!: string
	module!: !~ "license"
	action!: #StepAction
	cwd?: string
	params?: [...string]
} | {
	name!: string
	module!: =~ "license"
}

if #Step.module >= license {
}

#Steps: [...#Step]


#TemplateDef: {
  template: {
    filepath: string
    engine: string
  }
}
#Filename:=~ "^([a-zA-Z0-9_-]*\\.)+[a-zA-Z0-9_]+$"
#Complexfile:[#Filename]: #TemplateDef
#File: #Complexfile | #Filename

#Complexfolder:[string]: #FolderStruct
#Folder: #Complexfolder | string

#FolderStruct: [...#Folder|#File]


#Workflow: {
	config?: #Config
	vars?: #Vars
	steps?: #Steps
	folder_struct?: #FolderStruct
}

workflow: #Workflow
