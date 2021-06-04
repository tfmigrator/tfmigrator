package tfmigrator

import (
	"errors"

	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func ValidateResourceName(name string) error {
	if hclsyntax.ValidIdentifier(name) {
		return nil
	}
	return errors.New("invalid resource path: " + name)
}
