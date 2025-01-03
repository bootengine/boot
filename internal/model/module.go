package model

// a Module is the database representation of how third-party code is stored to be called by the application.
type Module struct {
	Name string     `db:"module_name"`
	Type ModuleType `db:"module_type"`
	Path string     `db:"module_path"`
}
