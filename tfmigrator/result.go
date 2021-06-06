package tfmigrator

import "path/filepath"

// Result contains a plan how resources are migrated.
// By marshaling Result as YAML, we can check the migration plan in advance.
type Result struct {
	Source           *Source
	MigratedResource *MigratedResource
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
