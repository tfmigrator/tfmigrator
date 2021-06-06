package tfmigrator

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

var validate = validator.New() //nolint:gochecknoglobals

// ValidateResourceName validates a Terraform Resource name.
func ValidateResourceName(name string) error {
	if hclsyntax.ValidIdentifier(name) {
		return nil
	}
	return errors.New("invalid resource name: " + name)
}
