package tfmigrator

type DryRunResult struct {
	MigratedResources []MigratedResource `yaml:"migrated_resources"`
	ExcludedResources []string           `yaml:"excluded_resources"`
	NoMatchResources  []string           `yaml:"no_match_resources"`
}

type MigratedResource struct {
	SourceResourcePath string `yaml:"source_resource_path"`
	DestResourcePath   string `yaml:"dest_resource_path"`
	TFBasename         string `yaml:"tf_basename"`
	StateDirname       string `yaml:"state_dirname"`
	StateBasename      string `yaml:"state_basename"`
}
