package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func (ctrl *Controller) Run(ctx context.Context, param Param) error {
	cfg := Config{}
	if err := ctrl.readConfig(param, &cfg); err != nil {
		return err
	}
	param.Items = cfg.Items
	state := State{}
	if param.StatePath != "" {
		if err := ctrl.readState(param.StatePath, &state); err != nil {
			return err
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

	for i, item := range param.Items {
		cr, err := ctrl.Matcher.Compile(item.Rule)
		if err != nil {
			return err
		}
		item.CompiledRule = cr
		if item.ResourceName != "" {
			crpc, err := ctrl.ResourcePathComputer.Compile(item.ResourceName)
			if err != nil {
				return err
			}
			item.CompiledResourceName = crpc
		}
		param.Items[i] = item
	}

	for _, rsc := range state.Values.RootModule.Resources {
		if err := ctrl.handleResource(ctx, param, rsc, tfPath); err != nil {
			return err
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

type ResourcePath struct {
	Type string
	Name string
}

func (rp *ResourcePath) Path() string {
	return rp.Type + "." + rp.Name
}

func (ctrl *Controller) handleResource(
	ctx context.Context, param Param, rsc Resource, hclFilePath string) error {
	for _, item := range param.Items {
		f, err := ctrl.handleItem(ctx, rsc, item, hclFilePath, param.SkipState)
		if err != nil {
			return err
		}
		if f {
			break
		}
	}
	return nil
}

func (ctrl *Controller) handleItem(
	ctx context.Context, rsc Resource, item Item, hclFilePath string, skipState bool) (bool, error) {
	// filter resource by condition
	matched, err := item.CompiledRule.Match(rsc)
	if err != nil {
		return false, err
	}
	if !matched {
		return false, nil
	}

	if item.Exclude {
		return true, nil
	}

	resourcePath, err := getResourcePath(rsc)
	if err != nil {
		return true, err
	}
	newResourcePath := resourcePath
	if item.ResourceName != "" {
		// compute new resource path
		newResourcePath.Name, err = item.CompiledResourceName.Parse(rsc)
		if err != nil {
			return true, err
		}
	}
	hclFile, err := os.Open(hclFilePath)
	if err != nil {
		return true, fmt.Errorf("open a Terraform configuration %s: %w", hclFilePath, err)
	}
	defer hclFile.Close()

	tfFile, err := os.OpenFile(item.TFPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return true, fmt.Errorf("open a file which will write Terraform configuration %s: %w", item.TFPath, err)
	}
	defer tfFile.Close()

	buf := bytes.Buffer{}
	if err := ctrl.getHCL(ctx, resourcePath.Path(), newResourcePath.Path(), hclFile, &buf); err != nil {
		return true, err
	}

	if err := ctrl.stateMv(ctx, item.StateOut, resourcePath.Path(), newResourcePath.Path(), skipState); err != nil {
		return true, err
	}
	// write hcl
	if _, err := io.Copy(tfFile, &buf); err != nil {
		return true, fmt.Errorf("write Terraform configuration to a file %s: %w", item.TFPath, err)
	}
	return true, nil
}

func (ctrl *Controller) getHCL(
	ctx context.Context, resourcePath, newResourcePath string, hclFile io.Reader, buf io.Writer) error {
	if resourcePath == newResourcePath {
		return ctrl.blockGet(ctx, "resource."+resourcePath, hclFile, buf)
	}
	pp := bytes.Buffer{}
	if err := ctrl.blockGet(ctx, "resource."+resourcePath, hclFile, &pp); err != nil {
		return err
	}

	if err := ctrl.blockMv(ctx, "resource."+resourcePath, "resource."+newResourcePath, &pp, buf); err != nil {
		return err
	}
	return nil
}

func getResourcePath(rsc Resource) (ResourcePath, error) { //nolint:unparam
	return ResourcePath{
		Type: rsc.Type,
		Name: rsc.Name,
	}, nil
}
