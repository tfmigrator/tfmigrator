package tfmigrator

import (
	"context"
	"fmt"

	"github.com/tfmigrator/tfmigrator/tfmigrator/hcledit"
	"github.com/tfmigrator/tfmigrator/tfmigrator/log"
	"github.com/tfmigrator/tfmigrator/tfmigrator/tfstate"
)

type BatchRunner struct {
	Planner     BatchPlanner `validate:"required"`
	Logger      log.Logger
	HCLEdit     *hcledit.Client `validate:"required"`
	StateReader *tfstate.Reader `validate:"required"`
	Migrator    *Migrator
	Outputter   Outputter
	DryRun      bool
}

// Run reads Terraform Configuration and State and migrate them.
func (runner *BatchRunner) Run(ctx context.Context, opt *RunOpt) error { //nolint:cyclop
	if err := validate.Struct(opt); err != nil {
		return fmt.Errorf("validate RunOpt: %w", err)
	}

	state, err := runner.StateReader.Read(ctx, opt.SourceStatePath)
	if err != nil {
		return fmt.Errorf("read Terraform State: %w", err)
	}

	addressFileMap, err := runner.HCLEdit.ListBlockMaps(opt.SourceHCLFilePaths...)
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

	results, err := runner.Planner.Plan(state, addressFileMap)
	if err != nil {
		return fmt.Errorf("plan: %w", err)
	}

	for _, result := range results {
		if result.Source.HCLFilePath == "" {
			result.Source.HCLFilePath = addressFileMap[result.Source.HCLAddress()]
		}
		if result.Source.StatePath == "" {
			result.Source.StatePath = opt.SourceStatePath
		}
		if err := runner.Migrator.Migrate(ctx, result.Source, result.MigratedResource); err != nil {
			return fmt.Errorf("migrate a resource: %w", err)
		}
	}

	if runner.Outputter != nil {
		if err := runner.Outputter.Output(results); err != nil {
			return fmt.Errorf("output the result: %w", err)
		}
	}

	return nil
}
