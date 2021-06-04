package tfmigrator

import (
	"errors"

	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// ValidateResourceName validates a Terraform Resource name.
func ValidateResourceName(name string) error {
	if hclsyntax.ValidIdentifier(name) {
		return nil
	}
	return errors.New("invalid resource name: " + name)
}
