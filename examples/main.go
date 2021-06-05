package main

import (
	"context"
	"log"
	"strings"

	"github.com/suzuki-shunsuke/tfmigrator-sdk/tfmigrator"
)

func main() {
	if err := core(); err != nil {
		log.Fatal(err)
	}
}

func core() error {
	ctx := context.Background()
	runner := &tfmigrator.Runner{}
	if err := runner.Run(ctx, &tfmigrator.RunOpt{
		Migrator: nil,
	}); err != nil {
		return err
	}
	return nil
}

var targets = []Matcher{
	NewStringMatcher("foo"),
}

type Matcher interface {
	StateDir(*tfmigrator.Resource) (string, error)
	Match(*tfmigrator.Resource) (bool, error)
}

type StringMatcher string

func NewStringMatcher(s string) *StringMatcher {
	a := StringMatcher(s)
	return &a
}

func (m *StringMatcher) StateDir(rsc *tfmigrator.Resource) (string, error) {
	return string(*m), nil
}

func (m *StringMatcher) Match(rsc *tfmigrator.Resource) (bool, error) {
	keys := strings.Split(string(*m), "/")
	for _, v := range rsc.Values {
		s, ok := v.(string)
		if !ok {
			continue
		}
		f := true
		for _, key := range keys {
			if !strings.Contains(s, key) {
				f = false
				break
			}
		}
		if f {
			return true, nil
		}
	}
	return false, nil
}

type Migrator struct{}

func (migrator *Migrator) Migrate(rsc *tfmigrator.Resource) (*tfmigrator.MigratedResource, error) {
	for _, target := range targets {
		f, err := target.Match(rsc)
		if err != nil {
			return nil, err
		}
		if !f {
			continue
		}
		stateDir, err := target.StateDir(rsc)
		if err != nil {
			return nil, err
		}
		if f {
			return &tfmigrator.MigratedResource{
				SourceAddress: rsc.Address,
				DestAddress:   getDestResourcePath(rsc),
				TFBasename:    rsc.Type + ".tf",
				StateDirname:  stateDir,
			}, nil
		}
	}
	return nil, nil
}

func getDestResourcePath(rsc *tfmigrator.Resource) string {
	switch rsc.Type {
	case "aws_iam_role":
		return rsc.Values["name"].(string)
	}
	return ""
}
