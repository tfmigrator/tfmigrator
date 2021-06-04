package controller

import (
	"errors"
	"fmt"
	"os"

	"github.com/suzuki-shunsuke/go-template-unmarshaler/text"
	"github.com/suzuki-shunsuke/tfmigrator/pkg/expr"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Items []Item
}

type Item struct {
	Rule          *expr.Bool
	Exclude       bool
	Stop          bool
	StateDirname  *text.Template `yaml:"state_dirname"`
	StateBasename *text.Template `yaml:"state_basename"`
	ResourceName  *ResourceName  `yaml:"resource_name"`
	TFBasename    *text.Template `yaml:"tf_basename"`
	Children      []Item
}

type MatchedItem struct {
	StateDirname  *text.Template
	StateBasename *text.Template
	ResourceName  *ResourceName
	TFBasename    *text.Template
	Exclude       bool
	Stop          bool
}

func (matchedItem *MatchedItem) Match() bool {
	return matchedItem.Exclude ||
		matchedItem.StateDirname != nil ||
		matchedItem.ResourceName != nil ||
		matchedItem.StateBasename != nil ||
		matchedItem.TFBasename != nil
}

func (matchedItem *MatchedItem) Parse(rsc Resource) (*MigratedResource, error) {
	resourceName := rsc.Name
	if !matchedItem.ResourceName.Empty() {
		name, err := matchedItem.ResourceName.Parse(rsc)
		if err != nil {
			return nil, err
		}
		resourceName = name
	}

	if matchedItem.TFBasename.Empty() {
		return nil, errors.New("tf_basename is required")
	}
	tfBasename, err := matchedItem.TFBasename.Execute(rsc)
	if err != nil {
		return nil, fmt.Errorf("render tf_basename: %w", err)
	}

	if matchedItem.StateDirname.Empty() {
		return nil, errors.New("state_dirname is required")
	}
	stateDirname, err := matchedItem.StateDirname.Execute(rsc)
	if err != nil {
		return nil, fmt.Errorf("render state_dirname: %w", err)
	}

	stateBasename := "terraform.tfstate"
	if matchedItem.StateBasename != nil {
		s, err := matchedItem.StateBasename.Execute(rsc)
		if err != nil {
			return nil, fmt.Errorf("render state_basename: %w", err)
		}
		stateBasename = s
	}

	return &MigratedResource{
		SourceResourcePath: rsc.Address,
		DestResourcePath:   rsc.Type + "." + resourceName,
		TFBasename:         tfBasename,
		StateDirname:       stateDirname,
		StateBasename:      stateBasename,
	}, nil
}

type Param struct {
	ConfigFilePath string
	LogLevel       string
	StatePath      string
	Items          []Item
	SkipState      bool
	DryRun         bool
}

type State struct {
	Values Values `json:"values"`
}

type Values struct {
	RootModule RootModule `json:"root_module"`
}

type RootModule struct {
	Resources []Resource `json:"resources"`
}

type Resource struct {
	Address string                 `json:"address"`
	Type    string                 `json:"type"`
	Name    string                 `json:"name"`
	Values  map[string]interface{} `json:"values"`
}

type DryRunResult struct {
	MigratedResources []MigratedResource `yaml:"migrated_resources"`
	ExcludedResources []string           `yaml:"excluded_resources"`
	NoMatchResources  []string           `yaml:"no_match_resources"`
}

type MigratedResource struct {
	SourceResourcePath string `yaml:"source_resource_path"`
	DestResourcePath   string `yaml:"dest_resource_path"`
	TFBasename         string `yaml:"tf_basename"`
	StateDirname       string `yaml:"state_dirname"`
	StateBasename      string `yaml:"state_basename"`
}

func (ctrl *Controller) readConfig(param Param, cfg *Config) error {
	cfgFile, err := os.Open(param.ConfigFilePath)
	if err != nil {
		return fmt.Errorf("open a configuration file %s: %w", param.ConfigFilePath, err)
	}
	defer cfgFile.Close()
	if err := yaml.NewDecoder(cfgFile).Decode(&cfg); err != nil {
		return fmt.Errorf("parse a configuration file as YAML %s: %w", param.ConfigFilePath, err)
	}
	return nil
}
