package tfmigrator

import (
	"context"
	"fmt"
	"io"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/tfmigrator/tfmigrator/tfmigrator/hcledit"
	"github.com/tfmigrator/tfmigrator/tfmigrator/log"
	"github.com/tfmigrator/tfmigrator/tfmigrator/tfstate"
)

type BulkRunner struct {
	Stdin        io.Reader   `validate:"required"`
	Stdout       io.Writer   `validate:"required"`
	Stderr       io.Writer   `validate:"required"`
	Planner      BulkPlanner `validate:"required"`
	Logger       log.Logger
	HCLEdit      *hcledit.Client  `validate:"required"`
	StateReader  *tfstate.Reader  `validate:"required"`
	StateUpdater *tfstate.Updater `validate:"required"`
	Reserver     *Reserver
	Outputter    Outputter
	DryRun       bool
}

type BulkPlanner interface {
	Plan(state *tfjson.State, addressFileMap map[string]string) error
}

type Reserver struct {
	results []Result
}

func (reserver *Reserver) Update(src *Source, migratedResource *MigratedResource) {
	reserver.results = append(reserver.results, Result{
		Source:           src,
		MigratedResource: migratedResource,
	})
}

// Run reads Terraform Configuration and State and migrate them.
func (runner *BulkRunner) Run(ctx context.Context, opt *RunOpt) error { //nolint:funlen,cyclop
	if err := validate.Struct(opt); err != nil {
		return fmt.Errorf("validate RunOpt: %w", err)
	}

	state, err := runner.StateReader.TFShow(ctx, opt.SourceStatePath)
	if err != nil {
		return fmt.Errorf("read Terraform State: %w", err)
	}

	addressFileMap, err := runner.HCLEdit.ListBlockMaps(opt.SourceTFFilePaths...)
	if err != nil {
		return fmt.Errorf("list all addresses in Terraform Configuration files: %w", err)
	}

	if state.Values == nil {
		runner.Logger.Info("state.Values is nil")
		return nil
	}
	if state.Values.RootModule == nil {
		runner.Logger.Info("state.Values.RootModule is nil")
		return nil
	}

	if err := runner.Planner.Plan(state, addressFileMap); err != nil {
		return err
	}

	for _, result := range runner.Reserver.results {
		if err := runner.Migrate(ctx, result.Source, result.MigratedResource); err != nil {
			return fmt.Errorf("migrate a resource: %w", err)
		}
	}

	if runner.Outputter != nil {
		if err := runner.Outputter.Output(runner.Reserver.results); err != nil {
			return fmt.Errorf("output the result: %w", err)
		}
	}

	return nil
}
