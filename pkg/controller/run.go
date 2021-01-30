package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

func (ctrl *Controller) Run(ctx context.Context, param Param) error {
	// read and validate parameter
	cfgFile, err := os.Open(param.ConfigFilePath)
	if err != nil {
		return err
	}
	defer cfgFile.Close()
	cfg := Config{}
	if err := yaml.NewDecoder(cfgFile).Decode(&cfg); err != nil {
		return err
	}
	param.Items = cfg.Items
	// read config
	// read resource from state
	// TODO compile rules in advance
	stateFile, err := os.Open(param.StatePath)
	if err != nil {
		return err
	}
	defer stateFile.Close()
	state := State{}
	if err := json.NewDecoder(stateFile).Decode(&state); err != nil {
		return err
	}

	f, err := ioutil.TempFile("", "")
	if err != nil {
		return err
	}
	defer f.Close()
	defer os.Remove(f.Name())
	// read tf from stdin and write a temporal file
	if _, err := io.Copy(f, ctrl.Stdin); err != nil {
		return err
	}

	for _, rsc := range state.Resources {
		if err := ctrl.handleResource(ctx, param, rsc, f.Name()); err != nil {
			return err
		}
	}
	return nil
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
	cr, err := ctrl.Matcher.Compile(item.Rule)
	if err != nil {
		return false, err
	}

	// filter resource by condition
	matched, err := cr.Match(rsc)
	if err != nil {
		return false, err
	}
	if !matched {
		return false, nil
	}
	resourcePath, err := getResourcePath(rsc)
	if err != nil {
		return true, err
	}
	newResourcePath := resourcePath
	if item.ResourceName != "" {
		// compute new resource path
		crpc, err := ctrl.ResourcePathComputer.Compile(item.ResourceName)
		if err != nil {
			return true, err
		}
		newResourcePath.Name, err = crpc.Parse(rsc)
		if err != nil {
			return true, err
		}
	}
	hclFile, err := os.Open(hclFilePath)
	if err != nil {
		return true, err
	}
	defer hclFile.Close()

	tfFile, err := os.OpenFile(item.TFPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return true, err
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
		return true, err
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

func getResourcePath(rsc Resource) (ResourcePath, error) {
	typ, ok := rsc["type"]
	if !ok {
		return ResourcePath{}, errors.New("state is invalid: resoruce type isn't found")
	}
	t, ok := typ.(string)
	if !ok {
		return ResourcePath{}, fmt.Errorf("state is invalid: resoruce type must be a string: %+v", typ)
	}

	name, ok := rsc["name"]
	if !ok {
		return ResourcePath{}, errors.New("state is invalid: resoruce name isn't found")
	}
	n, ok := name.(string)
	if !ok {
		return ResourcePath{}, fmt.Errorf("state is invalid: resoruce name must be a string: %+v", name)
	}
	return ResourcePath{
		Type: t,
		Name: n,
	}, nil
}
