package tfmigrator

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/tfmigrator/tfmigrator/tfmigrator/hcledit"
	"github.com/tfmigrator/tfmigrator/tfmigrator/log"
	"github.com/tfmigrator/tfmigrator/tfmigrator/tfstate"
)

// Runner provides high level API to migrate Terraform Configuration and State.
type Runner struct {
	Stdin        io.Reader `validate:"required"`
	Stdout       io.Writer `validate:"required"`
	Stderr       io.Writer `validate:"required"`
	Planner      Planner   `validate:"required"`
	Logger       log.Logger
	HCLEdit      *hcledit.Client  `validate:"required"`
	StateReader  *tfstate.Reader  `validate:"required"`
	StateUpdater *tfstate.Updater `validate:"required"`
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
			DryRun: runner.DryRun,
			Logger: runner.Logger,
		}
	}
	if runner.StateReader == nil {
		runner.StateReader = &tfstate.Reader{
			Stderr: runner.Stderr,
			Logger: runner.Logger,
		}
	}
	if runner.StateUpdater == nil {
		runner.StateUpdater = &tfstate.Updater{
			Stdout: runner.Stdout,
			Stderr: runner.Stderr,
			DryRun: runner.DryRun,
			Logger: runner.Logger,
		}
	}
}

// RunOpt is an option of Run method.
type RunOpt struct {
	// SourceStatePath is the file path to State.
	// If SourceStatePath is empty, State is read by `terraform show -json` command.
	SourceStatePath string
	// SourceTFFilePaths is a list of Terraform Configuration file paths.
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

	results := make([]Result, len(state.Values.RootModule.Resources))
	for i, rsc := range state.Values.RootModule.Resources {
		rsc := rsc
		tfFilePath := addressFileMap["resource."+rsc.Address]
		src := &Source{
			Resource:   &rsc,
			StatePath:  opt.SourceStatePath,
			TFFilePath: tfFilePath,
		}
		migratedResource, err := runner.migrateResource(ctx, src)
		if err != nil {
			return fmt.Errorf("migrate a resource %s: %w", rsc.Address, err)
		}
		if runner.Outputter != nil {
			results[i] = Result{
				Source:           src,
				MigratedResource: migratedResource,
			}
		}
	}

	if runner.Outputter != nil {
		if err := runner.Outputter.Output(results); err != nil {
			return fmt.Errorf("output the result: %w", err)
		}
	}

	return nil
}

func (runner *Runner) readState(ctx context.Context, sourceStatePath string, state *tfstate.State) error {
	if err := runner.StateReader.ReadByCmd(ctx, sourceStatePath, state); err != nil {
		return fmt.Errorf("read Terraform State by command: %w", err)
	}
	return nil
}

func (runner *Runner) migrateResource(ctx context.Context, source *Source) (*MigratedResource, error) {
	if err := validate.Struct(source); err != nil {
		return nil, fmt.Errorf("validate Source: %w", err)
	}
	migratedResource, err := runner.Planner.Plan(source)
	if err != nil {
		return nil, fmt.Errorf("plan to migrate a resource: %w", err)
	}

	if err := runner.Migrate(ctx, source, migratedResource); err != nil {
		return nil, fmt.Errorf("migrate a resource: %w", err)
	}
	return migratedResource, nil
}
