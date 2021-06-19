package tfmigrator

import (
	"path/filepath"
	"strings"
)

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
	// HCLFileBasename is a file name of Terraform Configuration.
	// StateBasename is a file name of Terraform State.
	// If Dirname and StateBasename is empty, the same State is updated.
	// If Dirname and HCLFileBasename is empty, the Terraform Configuration is updated in-place.
	Dirname         string
	HCLFileBasename string
	StateBasename   string
	// If Removed is true, the Resource is removed from Terraform Configuration and State.
	Removed            bool
	SkipHCLMigration   bool
	SkipStateMigration bool
}

// StatePath returns a file path to Terraform State file.
func (rsc *MigratedResource) StatePath() string {
	if rsc.Dirname != "" && rsc.StateBasename == "" {
		return filepath.Join(rsc.Dirname, "terraform.tfstate")
	}
	return filepath.Join(rsc.Dirname, rsc.StateBasename)
}

// HCLFilePath returns a file path to the Terraform Configuration file where the migrated Configuration is written.
func (rsc *MigratedResource) HCLFilePath() string {
	return filepath.Join(rsc.Dirname, rsc.HCLFileBasename)
}

// HCLAddress returns a resource address like `resource.null_resource.foo`.
func (rsc *MigratedResource) HCLAddress() string {
	if strings.HasPrefix(rsc.Address, "module.") {
		return rsc.Address
	}
	return "resource." + rsc.Address
}
