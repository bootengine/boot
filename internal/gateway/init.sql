CREATE TABLE  IF NOT EXISTS module (
	module_name TEXT PRIMARY KEY,
	module_path TEXT NOT NULL,
	module_type TEXT CHECK(module_type IN ("module", "filer", "vcs","template_engine")) NOT NULL
)

