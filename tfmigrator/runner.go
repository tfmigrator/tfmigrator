package tfmigrator

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

// Runner provides high level API to migrate Terraform Configuration and State.
type Runner struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
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

func (runner *Runner) dryRun(stdout io.Writer, dryRunResult *DryRunResult) error {
	if err := yaml.NewEncoder(stdout).Encode(dryRunResult); err != nil {
		return fmt.Errorf("output DryRunResult as YAML: %w", err)
	}
	return nil
}

// Run reads Terraform Configuration and State and migrate them.
func (runner *Runner) Run(ctx context.Context, opt *RunOpt) error {
	runner.SetDefault()
	if err := runner.validateOpt(opt); err != nil {
		return fmt.Errorf("validate a RunOpt: %w", err)
	}
	// read tf files from stdin
	tfFilePath, err := WriteTFInTemporalFile(runner.Stdin)
	if err != nil {
		return fmt.Errorf("write Terraform Configuration in a temporal file: %w", err)
	}
	defer os.Remove(tfFilePath)
	stdout := runner.Stdout
	stderr := runner.Stderr

	state := &State{}
	if opt.StatePath == "" {
		// read state by command
		if err := ReadStateByCmd(ctx, &ReadStateByCmdOpt{
			Stderr: stderr,
		}, state); err != nil {
			return fmt.Errorf("read Terraform State by command: %w", err)
		}
	} else {
		if err := ReadStateFromFile(opt.StatePath, state); err != nil {
			return fmt.Errorf("read Terraform State from a file %s: %w", opt.StatePath, err)
		}
	}

	dryRunResult := &DryRunResult{}
	for _, rsc := range state.Values.RootModule.Resources {
		rsc := rsc
		if err := runner.migrateResource(ctx, &rsc, &migrateResourceOpt{
			TFFilePath: tfFilePath,
			DryRun:     opt.DryRun,
			Migrator:   opt.Migrator,
		}, dryRunResult); err != nil {
			return fmt.Errorf("migrate a resource %s: %w", rsc.Address, err)
		}
	}

	if opt.DryRun {
		return runner.dryRun(stdout, dryRunResult)
	}

	return nil
}

type migrateResourceOpt struct {
	TFFilePath string
	Migrator   Migrator
	DryRun     bool
}

func (runner *Runner) migrateResource(ctx context.Context, rsc *Resource, opt *migrateResourceOpt, dryRunResult *DryRunResult) error {
	migratedResource, err := opt.Migrator.Migrate(rsc)
	if err != nil {
		return fmt.Errorf("plan to migrate a resource: %w", err)
	}
	if opt.DryRun {
		dryRunResult.Add(rsc.Address, migratedResource)
		return nil
	}

	if err := Migrate(ctx, migratedResource, &MigrateOpt{
		Stdin:      runner.Stdin,
		Stderr:     runner.Stderr,
		DryRun:     opt.DryRun,
		TFFilePath: opt.TFFilePath,
	}); err != nil {
		return fmt.Errorf("migrate a resource: %w", err)
	}
	return nil
}
