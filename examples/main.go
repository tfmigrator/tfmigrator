package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/tfmigrator-sdk/tfmigrator"
	"gopkg.in/yaml.v2"
)

func main() {
	if err := core(); err != nil {
		log.Fatal(err)
	}
}

func core() error {
	ctx := context.Background()
	dryRun := false
	logE := logrus.NewEntry(logrus.New())

	// read tf files from stdin
	tfFilePath, err := tfmigrator.WriteTFInTemporalFile(os.Stdin)
	if err != nil {
		return err
	}
	defer os.Remove(tfFilePath)
	stdin := os.Stdin
	stdout := os.Stdout
	stderr := os.Stderr

	// read state by command
	state := &tfmigrator.State{}
	if err := tfmigrator.ReadStateFromCmd(ctx, &tfmigrator.ReadStateFromCmdOpt{
		Stderr: stderr,
	}, state); err != nil {
		return err
	}

	dryRunResult := tfmigrator.DryRunResult{}
	for _, rsc := range state.Values.RootModule.Resources {
		migratedResource, err := migrateResource(&rsc)
		if err != nil {
			return err
		}
		if dryRun {
			dryRunResult.Add(migratedResource)
			continue
		}

		if err := tfmigrator.Migrate(ctx, migratedResource, &tfmigrator.MigrateOpt{
			Stdin:      stdin,
			Stderr:     stderr,
			DryRun:     dryRun,
			Logger:     logE,
			TFFilePath: tfFilePath,
		}); err != nil {
			return err
		}
	}

	if dryRun {
		if err := yaml.NewEncoder(stdout).Encode(dryRunResult); err != nil {
			return err
		}
		return nil
	}

	return nil
}

var targets = []Matcher{
	NewStringMatcher("heroku"),
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

type Target struct {
	Keys     []string
	StateDir string
}

func migrateResource(rsc *tfmigrator.Resource) (*tfmigrator.MigratedResource, error) {
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
				SourceResourcePath: rsc.Address,
				DestResourcePath:   getDestResourcePath(rsc),
				TFBasename:         rsc.Type + ".tf",
				StateDirname:       stateDir,
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
