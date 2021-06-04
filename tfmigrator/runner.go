package tfmigrator

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Runner provides high level API to migrate Terraform Configuration and State.
type Runner struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Logger *logrus.Entry
}

// SetDefault sets the default values to Runner.
func (runner *Runner) SetDefault() {
	if runner.Stdin == nil {
		runner.Stdin = os.Stdin
	}
	if runner.Stdout == nil {
		runner.Stdout = os.Stdout
	}
	if runner.Stderr == nil {
		runner.Stderr = os.Stderr
	}
	if runner.Stderr == nil {
		runner.Logger = logrus.NewEntry(logrus.New())
	}
}

var errMigratorIsRequired = errors.New("Migrator is required")

func (runner *Runner) validateOpt(opt *RunOpt) error {
	if opt.Migrator == nil {
		return errMigratorIsRequired
	}
	return nil
}

// RunOpt is an option of Run method.
type RunOpt struct {
	StatePath string
	DryRun    bool
	Migrator  Migrator
}

// Run reads Terraform Configuration and State and migrate them.
func (runner *Runner) Run(ctx context.Context, opt *RunOpt) error {
	runner.SetDefault()
	if err := runner.validateOpt(opt); err != nil {
		return err
	}
	// read tf files from stdin
	tfFilePath, err := WriteTFInTemporalFile(runner.Stdin)
	if err != nil {
		return fmt.Errorf("write Terraform Configuration in a temporal file: %w", err)
	}
	defer os.Remove(tfFilePath)
	stdin := runner.Stdin
	stdout := runner.Stdout
	stderr := runner.Stderr

	state := &State{}
	if opt.StatePath == "" {
		// read state by command
		if err := ReadStateByCmd(ctx, &ReadStateByCmdOpt{
			Stderr: stderr,
		}, state); err != nil {
			return err
		}
	} else {
		if err := ReadStateFromFile(opt.StatePath, state); err != nil {
			return err
		}
	}

	dryRunResult := DryRunResult{}
	for _, rsc := range state.Values.RootModule.Resources {
		migratedResource, err := opt.Migrator.Migrate(&rsc)
		if err != nil {
			return err
		}
		if opt.DryRun {
			dryRunResult.Add(rsc.Address, migratedResource)
			continue
		}

		if err := Migrate(ctx, migratedResource, &MigrateOpt{
			Stdin:      stdin,
			Stderr:     stderr,
			DryRun:     opt.DryRun,
			Logger:     runner.Logger,
			TFFilePath: tfFilePath,
		}); err != nil {
			return err
		}
	}

	if opt.DryRun {
		if err := yaml.NewEncoder(stdout).Encode(dryRunResult); err != nil {
			return err
		}
		return nil
	}

	return nil
}
