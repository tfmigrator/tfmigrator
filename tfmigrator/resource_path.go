package controller

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type ResourceName struct {
	tmpl *template.Template
	raw  string
}

func (resourceName *ResourceName) Empty() bool {
	return resourceName == nil || resourceName.tmpl == nil
}

func (resourceName *ResourceName) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var a string
	if err := unmarshal(&a); err != nil {
		return err
	}
	t, err := template.New("_").Funcs(sprig.TxtFuncMap()).Parse(a)
	if err != nil {
		return fmt.Errorf("parse a string with text/template: %w", err)
	}
	resourceName.raw = a
	resourceName.tmpl = t
	return nil
}

func (resourceName *ResourceName) Parse(rsc interface{}) (string, error) {
	buf := &bytes.Buffer{}
	if err := resourceName.tmpl.Execute(buf, rsc); err != nil {
		return "", fmt.Errorf("render a template with params: %w", err)
	}
	p := buf.String()
	if !hclsyntax.ValidIdentifier(p) {
		return "", fmt.Errorf("invalid resource path: " + p)
	}
	return buf.String(), nil
}

func (resourceName *ResourceName) Raw() string {
	return resourceName.raw
}
