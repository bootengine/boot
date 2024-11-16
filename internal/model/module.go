package model

type Module struct {
	Name string     `db:"module_name"`
	Type ModuleType `db:"module_type"`
	Path string     `db:"module_path"`
}
