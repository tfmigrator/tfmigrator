package tfmigrator

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

// Outputter outputs the result.
type Outputter interface {
	Output([]Result) error
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
func (outputter *YAMLOutputter) Output(results []Result) error {
	if err := yaml.NewEncoder(outputter.out).Encode(outputter.format(results)); err != nil {
		return fmt.Errorf("output Result as YAML: %w", err)
	}
	return nil
}

func (outputter *YAMLOutputter) format(results []Result) *yamlResults {
	yr := &yamlResults{}
	for _, result := range results {
		rsc := result.MigratedResource
		src := result.Source
		if rsc == nil {
			a := yamlNotMigratedResult{
				Address:  src.Address(),
				FilePath: src.TFFilePath,
			}
			if src.Resource != nil {
				a.Attributes = src.Resource.AttributeValues
			}
			yr.NotMigratedResources = append(yr.NotMigratedResources, a)
			continue
		}
		if rsc.Removed {
			yr.RemovedResources = append(yr.RemovedResources, yamlSourceResult{
				Address:  src.Address(),
				FilePath: src.TFFilePath,
			})
			continue
		}
		yr.MigratedResources = append(yr.MigratedResources, yamlResult{
			SourceAddress:     src.Address(),
			SourceTFFilePath:  src.TFFilePath,
			NewAddress:        rsc.Address,
			NewTFFileBasename: rsc.TFFileBasename,
			Dirname:           rsc.Dirname,
			StateBasename:     rsc.StateBasename,
		})
	}
	return yr
}

type yamlResults struct {
	MigratedResources    []yamlResult            `yaml:"migrated_resources"`
	RemovedResources     []yamlSourceResult      `yaml:"removed_resources"`
	NotMigratedResources []yamlNotMigratedResult `yaml:"not_migrated_resources"`
}

type yamlResult struct {
	SourceAddress     string `yaml:"source_address"`
	SourceTFFilePath  string `yaml:"source_tf_file_path,omitempty"`
	NewAddress        string `yaml:"new_address,omitempty"`
	NewTFFileBasename string `yaml:"new_tf_file_basename,omitempty"`
	Dirname           string `yaml:"dirname,omitempty"`
	StateBasename     string `yaml:"state_basename,omitempty"`
}

type yamlSourceResult struct {
	Address  string
	FilePath string `yaml:"file_path,omitempty"`
}

type yamlNotMigratedResult struct {
	Address    string
	FilePath   string                 `yaml:"file_path,omitempty"`
	Attributes map[string]interface{} `yaml:",omitempty"`
}
