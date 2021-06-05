package tfmigrator

import (
	"context"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

// Runner provides high level API to migrate Terraform Configuration and State.
type Runner struct {
	Stdin    io.Reader `validate:"required"`
	Stdout   io.Writer `validate:"required"`
	Stderr   io.Writer `validate:"required"`
	Migrator Migrator  `validate:"required"`
	Logger   Logger
	DryRun   bool
}

func (runner *Runner) Validate() error {
	runner.SetDefault()
	if err := validate.Struct(runner); err != nil {
		return fmt.Errorf("validate Runner: %w", err)
	}
	return nil
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

// RunOpt is an option of Run method.
type RunOpt struct {
	SourceStatePath   string
	SourceTFFilePaths []string `validate:"required"`
}

func (runner *Runner) dryRun(stdout io.Writer, dryRunResult *DryRunResult) error {
	if err := yaml.NewEncoder(stdout).Encode(dryRunResult); err != nil {
		return fmt.Errorf("output DryRunResult as YAML: %w", err)
	}
	return nil
}

// Run reads Terraform Configuration and State and migrate them.
func (runner *Runner) Run(ctx context.Context, opt *RunOpt) error {
	if err := validate.Struct(opt); err != nil {
		return fmt.Errorf("validate RunOpt: %w", err)
	}
	runner.SetDefault()
	stdout := runner.Stdout
	stderr := runner.Stderr

	state := &State{}
	if opt.SourceStatePath == "" {
		// read state by command
		if err := ReadStateByCmd(ctx, &ReadStateByCmdOpt{
			Stderr: stderr,
		}, state); err != nil {
			return fmt.Errorf("read Terraform State by command: %w", err)
		}
	} else {
		if err := ReadStateFromFile(opt.SourceStatePath, state); err != nil {
			return fmt.Errorf("read Terraform State from a file %s: %w", opt.SourceStatePath, err)
		}
	}

	addressFileMap, err := listBlockMaps(&listBlockMapsOpt{
		FilePaths: opt.SourceTFFilePaths,
		Stderr:    stderr,
	})
	if err != nil {
		return err
	}

	dryRunResult := &DryRunResult{}
	for _, rsc := range state.Values.RootModule.Resources {
		rsc := rsc
		tfFilePath, ok := addressFileMap["resource."+rsc.Address]
		if !ok {
			continue
		}
		if err := runner.migrateResource(ctx, &Source{
			Resource:   &rsc,
			StatePath:  opt.SourceStatePath,
			TFFilePath: tfFilePath,
		}, dryRunResult); err != nil {
			return fmt.Errorf("migrate a resource %s: %w", rsc.Address, err)
		}
	}

	if runner.DryRun {
		return runner.dryRun(stdout, dryRunResult)
	}

	return nil
}

type Source struct {
	Resource *Resource
	// If the resource isn't found in Terraform Configuration files, TFFilePath is empty
	TFFilePath string
	StatePath  string
}

func (src *Source) Address() string {
	return src.Resource.Address
}

func (runner *Runner) migrateResource(ctx context.Context, source *Source, dryRunResult *DryRunResult) error {
	if err := validate.Struct(source); err != nil {
		return fmt.Errorf("validate Source: %w", err)
	}
	migratedResource, err := runner.Migrator.Migrate(source)
	if err != nil {
		return fmt.Errorf("plan to migrate a resource: %w", err)
	}
	if runner.DryRun {
		dryRunResult.Add(source, migratedResource)
		return nil
	}

	if err := runner.Migrate(ctx, source, migratedResource); err != nil {
		return fmt.Errorf("migrate a resource: %w", err)
	}
	return nil
}
