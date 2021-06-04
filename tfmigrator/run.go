package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func (ctrl *Controller) Run(ctx context.Context, param Param) error { //nolint:cyclop
	cfg := Config{}
	if err := ctrl.readConfig(param, &cfg); err != nil {
		return err
	}
	param.Items = cfg.Items
	state := State{}
	if param.StatePath != "" {
		if err := ctrl.readState(param.StatePath, &state); err != nil {
			return fmt.Errorf("read state (state path: %s): %w", param.StatePath, err)
		}
	} else {
		if err := ctrl.readStateFromCmd(ctx, &state); err != nil {
			return err
		}
	}

	tfPath, err := ctrl.writeTF()
	if tfPath != "" {
		defer os.Remove(tfPath)
	}
	if err != nil {
		return err
	}

	dryRunResult := DryRunResult{}

	for _, rsc := range state.Values.RootModule.Resources {
		if err := ctrl.handleResource(ctx, param, rsc, tfPath, &dryRunResult); err != nil {
			return fmt.Errorf("handle a resource %s: %w", rsc.Address, err)
		}
	}
	if param.DryRun {
		if err := yaml.NewEncoder(ctrl.Stdout).Encode(dryRunResult); err != nil {
			return fmt.Errorf("encode dry run result as YAML: %w", err)
		}
	}
	return nil
}

func (ctrl *Controller) readState(statePath string, state *State) error {
	stateFile, err := os.Open(statePath)
	if err != nil {
		return fmt.Errorf("open a state file %s: %w", statePath, err)
	}
	defer stateFile.Close()
	if err := json.NewDecoder(stateFile).Decode(state); err != nil {
		return fmt.Errorf("parse a state file as JSON %s: %w", statePath, err)
	}
	return nil
}

func (ctrl *Controller) readStateFromCmd(ctx context.Context, state *State) error {
	buf := bytes.Buffer{}
	if err := ctrl.tfShow(ctx, &buf); err != nil {
		return err
	}
	if err := json.NewDecoder(&buf).Decode(state); err != nil {
		return fmt.Errorf("parse a state as JSON: %w", err)
	}
	return nil
}

func (ctrl *Controller) writeTF() (string, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return "", fmt.Errorf("create a temporal file to write Terraform configuration (.tf): %w", err)
	}
	defer f.Close()
	// read tf from stdin and write a temporal file
	if _, err := io.Copy(f, ctrl.Stdin); err != nil {
		return f.Name(), fmt.Errorf("write standard input to a temporal file %s: %w", f.Name(), err)
	}
	return f.Name(), nil
}

func (ctrl *Controller) handleResource(ctx context.Context, param Param, rsc Resource, hclFilePath string, dryRunResult *DryRunResult) error { //nolint:funlen,cyclop
	matchedItem := MatchedItem{}
	for _, item := range param.Items {
		if err := ctrl.handleItem(rsc, item, &matchedItem); err != nil {
			return fmt.Errorf("handle item (rule: %s): %w", item.Rule.Raw(), err)
		}
		if matchedItem.Exclude || matchedItem.Stop {
			break
		}
	}

	if param.DryRun {
		if matchedItem.Exclude {
			dryRunResult.ExcludedResources = append(dryRunResult.ExcludedResources, rsc.Address)
			return nil
		}
		if !matchedItem.Match() {
			dryRunResult.NoMatchResources = append(dryRunResult.NoMatchResources, rsc.Address)
			return nil
		}
		migratedResource, err := matchedItem.Parse(rsc)
		if err != nil {
			return err
		}
		dryRunResult.MigratedResources = append(dryRunResult.MigratedResources, *migratedResource)
		return nil
	}

	if matchedItem.Exclude {
		return nil
	}

	if !matchedItem.Match() {
		return nil
	}
	migratedResource, err := matchedItem.Parse(rsc)
	if err != nil {
		return err
	}

	tfPath := filepath.Join(migratedResource.StateDirname, migratedResource.TFBasename)
	tfFile, err := os.OpenFile(tfPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open a file which will write Terraform configuration %s: %w", tfPath, err)
	}
	defer tfFile.Close()

	hclFile, err := os.Open(hclFilePath)
	if err != nil {
		return fmt.Errorf("open a Terraform configuration %s: %w", hclFilePath, err)
	}
	defer hclFile.Close()

	buf := bytes.Buffer{}
	if err := ctrl.getHCL(ctx, rsc.Address, migratedResource.DestResourcePath, hclFile, &buf); err != nil {
		return err
	}

	if err := ctrl.stateMv(ctx, filepath.Join(migratedResource.StateDirname, migratedResource.StateBasename), rsc.Address, migratedResource.DestResourcePath, param.SkipState); err != nil {
		return err
	}
	// write hcl
	if _, err := io.Copy(tfFile, &buf); err != nil {
		return fmt.Errorf("write Terraform configuration to a file %s: %w", tfPath, err)
	}

	return nil
}

func (ctrl *Controller) handleItem(rsc Resource, item Item, matchedItem *MatchedItem) error { //nolint:cyclop
	// filter resource by condition
	matched, err := item.Rule.Run(rsc)
	if err != nil {
		return fmt.Errorf("check if the rule matches with the resource: %w", err)
	}
	if !matched {
		return nil
	}

	if item.Exclude {
		matchedItem.Exclude = true
		return nil
	}

	if item.ResourceName != nil {
		matchedItem.ResourceName = item.ResourceName
	}
	if item.StateBasename != nil {
		matchedItem.StateBasename = item.StateBasename
	}
	if item.StateDirname != nil {
		matchedItem.StateDirname = item.StateDirname
	}
	if item.TFBasename != nil {
		matchedItem.TFBasename = item.TFBasename
	}

	if item.Stop {
		matchedItem.Stop = true
		return nil
	}

	for _, child := range item.Children {
		if err := ctrl.handleItem(rsc, child, matchedItem); err != nil {
			return err
		}
		if matchedItem.Stop || matchedItem.Exclude {
			return nil
		}
	}
	return nil
}

func (ctrl *Controller) getHCL(
	ctx context.Context, resourcePath, newResourcePath string, hclFile io.Reader, buf io.Writer) error {
	if resourcePath == newResourcePath {
		return ctrl.blockGet(ctx, "resource."+resourcePath, hclFile, buf)
	}
	pp := bytes.Buffer{}
	if err := ctrl.blockGet(ctx, "resource."+resourcePath, hclFile, &pp); err != nil {
		return fmt.Errorf("get a resource from HCL file: %w", err)
	}

	if err := ctrl.blockMv(ctx, "resource."+resourcePath, "resource."+newResourcePath, &pp, buf); err != nil {
		return fmt.Errorf("rename resource: %w", err)
	}
	return nil
}
