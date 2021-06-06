package tfmigrator

import "path/filepath"

// Result contains source Terraform Resource and the plan how the resource is migrated.
// Result is used to output the result by Outputter.
type Result struct {
	Source           *Source
	MigratedResource *MigratedResource
}

// MigratedResource is a plan how a resource is migrated.
type MigratedResource struct {
	// Address is a new resource address.
	// If Address is empty, the address isn't changed.
	Address string
	// Dirname is a directory path where Terraform Configuration file and State exists.
	// TFFileBasename is a file name of Terraform Configuration.
	// StateBasename is a file name of Terraform State.
	// If Dirname and StateBasename is empty, the same State is updated.
	// If Dirname and TFFileBasename is empty, the Terraform Configuration is updated in-place.
	Dirname        string
	TFFileBasename string
	StateBasename  string
	// If Removed is true, the Resource is removed from Terraform Configuration and State.
	Removed bool
}

// StatePath returns a file path to Terraform State file.
func (rsc *MigratedResource) StatePath() string {
	return filepath.Join(rsc.Dirname, rsc.StateBasename)
}

// TFFilePath returns a file path to the Terraform Configuration file where the migrated Configuration is written.
func (rsc *MigratedResource) TFFilePath() string {
	return filepath.Join(rsc.Dirname, rsc.TFFileBasename)
}
