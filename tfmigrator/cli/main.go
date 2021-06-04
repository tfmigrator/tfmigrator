package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/tfmigrator-sdk/tfmigrator"
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

// Migrator migrates a Terraform resource.
// Note that Migrator doesn't change Terraform State and Terraform Configuration files.
// Migrator determines the updated resource name, outputted State file path, and outputted Terraform Configuration file path.
// If migrator
type Migrator interface {
	Migrate(rsc *tfmigrator.Resource) (*tfmigrator.MigratedResource, error)
}

type combinedMigrator struct {
	migrators []Migrator
}

func CombineMigrators(migrators ...Migrator) Migrator {
	return &combinedMigrator{
		migrators: migrators,
	}
}

func (migrator *combinedMigrator) Migrate(rsc *tfmigrator.Resource) (*tfmigrator.MigratedResource, error) {
	for _, m := range migrator.migrators {
		migratedResource, err := m.Migrate(rsc)
		if err != nil {
			return nil, err
		}
		if migratedResource == nil {
			continue
		}
		return migratedResource, nil
	}
	return nil, nil
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
	tfFilePath, err := tfmigrator.WriteTFInTemporalFile(runner.Stdin)
	if err != nil {
		return fmt.Errorf("write Terraform Configuration in a temporal file: %w", err)
	}
	defer os.Remove(tfFilePath)
	stdin := runner.Stdin
	stdout := runner.Stdout
	stderr := runner.Stderr

	state := &tfmigrator.State{}
	if opt.StatePath == "" {
		// read state by command
		if err := tfmigrator.ReadStateFromCmd(ctx, &tfmigrator.ReadStateFromCmdOpt{
			Stderr: stderr,
		}, state); err != nil {
			return err
		}
	} else {
		if err := tfmigrator.ReadStateFromFile(opt.StatePath, state); err != nil {
			return err
		}
	}

	dryRunResult := tfmigrator.DryRunResult{}
	for _, rsc := range state.Values.RootModule.Resources {
		migratedResource, err := opt.Migrator.Migrate(&rsc)
		if err != nil {
			return err
		}
		if opt.DryRun {
			dryRunResult.Add(rsc.Address, migratedResource)
			continue
		}

		if err := tfmigrator.Migrate(ctx, migratedResource, &tfmigrator.MigrateOpt{
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
