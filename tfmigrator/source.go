package tfmigrator

import "github.com/suzuki-shunsuke/tfmigrator-sdk/tfmigrator/tfstate"

// Source is a source Terraform resource.
type Source struct {
	Resource *tfstate.Resource
	// If the resource isn't found in Terraform Configuration files, TFFilePath is empty
	TFFilePath string
	StatePath  string
}

// Address returns a resource address like `null_resource.foo`.
func (src *Source) Address() string {
	return src.Resource.Address
}
