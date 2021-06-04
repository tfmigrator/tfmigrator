package tfmigrator

import "path/filepath"

// DryRunResult contains a plan how resources are migrated.
// By marshaling DryRunResult as YAML, we can check the migration plan in advance.
type DryRunResult struct {
	MigratedResources []MigratedResource `yaml:"migrated_resources"`
	ExcludedResources []string           `yaml:"excluded_resources"`
	NoMatchResources  []string           `yaml:"no_match_resources"`
}

// Add adds a migration plan of a resource to DryRunResult.
func (result *DryRunResult) Add(rsc *MigratedResource) {
	if rsc.Exclude {
		result.ExcludedResources = append(result.ExcludedResources, rsc.SourceResourcePath)
		return
	}
	if rsc.NoMatch {
		result.NoMatchResources = append(result.NoMatchResources, rsc.SourceResourcePath)
		return
	}
	result.MigratedResources = append(result.MigratedResources, *rsc)
}

// MigratedResource is a plan how a resource is migrated
type MigratedResource struct {
	SourceResourcePath string `yaml:"source_resource_path"`
	DestResourcePath   string `yaml:"dest_resource_path"`
	TFBasename         string `yaml:"tf_basename"`
	StateDirname       string `yaml:"state_dirname"`
	StateBasename      string `yaml:"state_basename"`
	Exclude            bool   `yaml:"-"`
	NoMatch            bool   `yaml:"-"`
}

// StatePath returns a file path to Terraform State file.
func (rsc *MigratedResource) StatePath() string {
	return filepath.Join(rsc.StateDirname, rsc.StateBasename)
}

// TFPath returns a file path to the Terraform Configuration file where the migrated Configuration is written.
func (rsc *MigratedResource) TFPath() string {
	return filepath.Join(rsc.StateDirname, rsc.TFBasename)
}

// PathChanged returns true if the resource path is changed.
func (rsc *MigratedResource) PathChanged() bool {
	return rsc.SourceResourcePath != rsc.DestResourcePath
}
