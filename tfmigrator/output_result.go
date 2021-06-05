package tfmigrator

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

// Outputter outputs the result.
type Outputter interface {
	Output(*Result) error
}

// YAMLOutputter outputs the result with YAML format.
type YAMLOutputter struct {
	out io.Writer
}

// NewYAMLOutputter returns a YAMLOutputter.
func NewYAMLOutputter(out io.Writer) Outputter {
	return &YAMLOutputter{
		out: out,
	}
}

// Output outputs the result with YAML format.
func (outputter *YAMLOutputter) Output(result *Result) error {
	if err := yaml.NewEncoder(outputter.out).Encode(result); err != nil {
		return fmt.Errorf("output Result as YAML: %w", err)
	}
	return nil
}
