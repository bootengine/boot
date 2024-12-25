#Config : {
	create_root?: bool | true
	unrestricted?: bool | false
}

#Var : {
	type: "string" | "license" | "password"
	required: bool
}

#Vars: [string]:{
	#Var
}


#StepAction: "init" | "install-local-deps" | "install-global-deps" | "install-dev-deps" | "commit"| "push"| "add"| "add-origin"| "create-file"| "create-folder"| "write-file"| "apply-template" | "create-folder-struct"

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
	from?: string
}

workflow: #Workflow
