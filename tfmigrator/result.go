package tfmigrator

import "path/filepath"

// Result contains a plan how resources are migrated.
// By marshaling Result as YAML, we can check the migration plan in advance.
type Result struct {
	MigratedResources    []DryRunResource `yaml:"migrated_resources"`
	RemovedResources     []SourceResource `yaml:"removed_resources"`
	NotMigratedResources []SourceResource `yaml:"not_migrated_resources"`
}

type SourceResource struct {
	Address  string
	FilePath string `yaml:"file_path,omitempty"`
}

// Add adds a migration plan of a resource to Result.
func (result *Result) Add(src *Source, rsc *MigratedResource) {
	if rsc == nil {
		result.NotMigratedResources = append(result.NotMigratedResources, SourceResource{
			Address:  src.Address(),
			FilePath: src.TFFilePath,
		})
		return
	}
	if rsc.Removed {
		result.RemovedResources = append(result.RemovedResources, SourceResource{
			Address:  src.Address(),
			FilePath: src.TFFilePath,
		})
	}
	result.MigratedResources = append(result.MigratedResources, DryRunResource{
		SourceAddress:     src.Address(),
		SourceTFFilePath:  src.TFFilePath,
		NewAddress:        rsc.Address,
		NewTFFileBasename: rsc.TFFileBasename,
		Dirname:           rsc.Dirname,
		StateBasename:     rsc.StateBasename,
	})
}

// DryRunResource is a plan how a resource is migrated
type DryRunResource struct {
	SourceAddress     string `yaml:"source_address"`
	SourceTFFilePath  string `yaml:"source_tf_file_path,omitempty"`
	NewAddress        string `yaml:"new_address,omitempty"`
	NewTFFileBasename string `yaml:"new_tf_file_basename,omitempty"`
	Dirname           string `yaml:"dirname,omitempty"`
	StateBasename     string `yaml:"state_basename,omitempty"`
}

// MigratedResource is a plan how a resource is migrated
type MigratedResource struct {
	Address        string
	Dirname        string
	TFFileBasename string
	StateBasename  string
	Removed        bool
}

// StatePath returns a file path to Terraform State file.
func (rsc *MigratedResource) StatePath() string {
	return filepath.Join(rsc.Dirname, rsc.StateBasename)
}

// TFFilePath returns a file path to the Terraform Configuration file where the migrated Configuration is written.
func (rsc *MigratedResource) TFFilePath() string {
	return filepath.Join(rsc.Dirname, rsc.TFFileBasename)
}
