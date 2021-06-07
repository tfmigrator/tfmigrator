package tfmigrator

import tfjson "github.com/hashicorp/terraform-json"

// Source is a source Terraform resource.
type Source struct {
	Resource *tfjson.StateResource
	// If the resource isn't found in Terraform Configuration files, TFFilePath is empty
	TFFilePath string
	StatePath  string
}

// Address returns a resource address like `null_resource.foo`.
func (src *Source) Address() string {
	return src.Resource.Address
}
