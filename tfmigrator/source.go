package tfmigrator

import tfjson "github.com/hashicorp/terraform-json"

// Source is a source Terraform resource.
type Source struct {
	Resource *tfjson.StateResource
	Module   *tfjson.StateModule
	// If the resource isn't found in Terraform Configuration files, TFFilePath is empty
	TFFilePath string
	StatePath  string
}

// Address returns a resource address like `null_resource.foo`.
func (src *Source) Address() string {
	if src.Resource != nil {
		return src.Resource.Address
	}
	if src.Module != nil {
		return src.Module.Address
	}
	return ""
}

// HCLAddress returns a resource address like `resource.null_resource.foo`.
func (src *Source) HCLAddress() string {
	if src.Resource != nil {
		return "resource." + src.Resource.Address
	}
	if src.Module != nil {
		return src.Module.Address
	}
	return ""
}
