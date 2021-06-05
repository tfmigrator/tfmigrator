package tfmigrator

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/suzuki-shunsuke/tfmigrator-sdk/tfmigrator/hcledit"
	"github.com/suzuki-shunsuke/tfmigrator-sdk/tfmigrator/log"
	"github.com/suzuki-shunsuke/tfmigrator-sdk/tfmigrator/tfstate"
)

// Runner provides high level API to migrate Terraform Configuration and State.
type Runner struct {
	Stdin        io.Reader `validate:"required"`
	Stdout       io.Writer `validate:"required"`
	Stderr       io.Writer `validate:"required"`
	Planner      Planner   `validate:"required"`
	Logger       log.Logger
	HCLEdit      *hcledit.Client
	StateReader  *tfstate.Reader
	StateUpdater *tfstate.Updater
	Outputter    Outputter
	DryRun       bool
}

// Validate sets default values and validates runner.
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
	if runner.HCLEdit == nil {
		runner.HCLEdit = &hcledit.Client{
			Stderr: runner.Stderr,
		}
	}
	if runner.StateReader == nil {
		runner.StateReader = &tfstate.Reader{
			Stderr: runner.Stderr,
		}
	}
}

// RunOpt is an option of Run method.
type RunOpt struct {
	SourceStatePath   string
	SourceTFFilePaths []string `validate:"required"`
}

// Run reads Terraform Configuration and State and migrate them.
func (runner *Runner) Run(ctx context.Context, opt *RunOpt) error {
	if err := validate.Struct(opt); err != nil {
		return fmt.Errorf("validate RunOpt: %w", err)
	}
	runner.SetDefault()

	state := &tfstate.State{}
	if err := runner.readState(ctx, opt.SourceStatePath, state); err != nil {
		return err
	}

	addressFileMap, err := runner.HCLEdit.ListBlockMaps(opt.SourceTFFilePaths...)
	if err != nil {
		return fmt.Errorf("list all addresses in Terraform Configuration files: %w", err)
	}

	dryRunResult := &Result{}
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

	if runner.Outputter != nil {
		if err := runner.Outputter.Output(dryRunResult); err != nil {
			return fmt.Errorf("output the result: %w", err)
		}
	}

	return nil
}

func (runner *Runner) readState(ctx context.Context, sourceStatePath string, state *tfstate.State) error {
	if sourceStatePath == "" {
		// read state by command
		if err := runner.StateReader.ReadByCmd(ctx, state); err != nil {
			return fmt.Errorf("read Terraform State by command: %w", err)
		}
	} else {
		if err := tfstate.ReadFromFile(sourceStatePath, state); err != nil {
			return fmt.Errorf("read Terraform State from a file %s: %w", sourceStatePath, err)
		}
	}
	return nil
}

func (runner *Runner) migrateResource(ctx context.Context, source *Source, dryRunResult *Result) error {
	if err := validate.Struct(source); err != nil {
		return fmt.Errorf("validate Source: %w", err)
	}
	migratedResource, err := runner.Planner.Plan(source)
	if err != nil {
		return fmt.Errorf("plan to migrate a resource: %w", err)
	}
	if runner.Outputter != nil {
		dryRunResult.Add(source, migratedResource)
		return nil
	}

	if err := runner.Migrate(ctx, source, migratedResource); err != nil {
		return fmt.Errorf("migrate a resource: %w", err)
	}
	return nil
}
