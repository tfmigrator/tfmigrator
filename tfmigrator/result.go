package tfmigrator

import "path/filepath"

type DryRunResult struct {
	MigratedResources []MigratedResource `yaml:"migrated_resources"`
	ExcludedResources []string           `yaml:"excluded_resources"`
	NoMatchResources  []string           `yaml:"no_match_resources"`
}

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

type MigratedResource struct {
	SourceResourcePath string `yaml:"source_resource_path"`
	DestResourcePath   string `yaml:"dest_resource_path"`
	TFBasename         string `yaml:"tf_basename"`
	StateDirname       string `yaml:"state_dirname"`
	StateBasename      string `yaml:"state_basename"`
	Exclude            bool   `yaml:"-"`
	NoMatch            bool   `yaml:"-"`
}

func (rsc *MigratedResource) StatePath() string {
	return filepath.Join(rsc.StateDirname, rsc.StateBasename)
}

func (rsc *MigratedResource) TFPath() string {
	return filepath.Join(rsc.StateDirname, rsc.TFBasename)
}

func (rsc *MigratedResource) PathChanged() bool {
	return rsc.SourceResourcePath != rsc.DestResourcePath
}
