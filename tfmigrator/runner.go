package tfmigrator

import (
	"context"
	"fmt"

	"github.com/tfmigrator/tfmigrator/tfmigrator/hcledit"
	"github.com/tfmigrator/tfmigrator/tfmigrator/log"
	"github.com/tfmigrator/tfmigrator/tfmigrator/tfstate"
)

// Runner provides high level API to migrate Terraform Configuration and State.
type Runner struct {
	Planner     Planner `validate:"required"`
	Logger      log.Logger
	HCLEdit     *hcledit.Client `validate:"required"`
	StateReader *tfstate.Reader `validate:"required"`
	Outputter   Outputter
	Migrator    *Migrator
	DryRun      bool
}

// RunOpt is an option of Runner#Run method.
type RunOpt struct {
	// SourceStatePath is the file path to State.
	// If SourceStatePath is empty, State is read by `terraform show -json` command.
	SourceStatePath string
	// SourceHCLFilePaths is a list of Terraform Configuration file paths.
	SourceHCLFilePaths []string `validate:"required"`
}

// Run reads Terraform Configuration and State and migrate them.
func (runner *Runner) Run(ctx context.Context, opt *RunOpt) error { //nolint:funlen,cyclop
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

	results := make([]Result, 0, len(state.Values.RootModule.Resources)+len(state.Values.RootModule.ChildModules))
	for _, rsc := range state.Values.RootModule.Resources {
		tfFilePath := addressFileMap["resource."+rsc.Address]
		src := &Source{
			Resource:    rsc,
			StatePath:   opt.SourceStatePath,
			HCLFilePath: tfFilePath,
		}
		migratedResource, err := runner.migrateResource(ctx, src)
		if err != nil {
			return fmt.Errorf("migrate a resource %s: %w", rsc.Address, err)
		}
		if runner.Outputter != nil {
			results = append(results, Result{
				Source:           src,
				MigratedResource: migratedResource,
			})
		}
	}

	// update module addresses
	for _, child := range state.Values.RootModule.ChildModules {
		tfFilePath := addressFileMap[child.Address]
		src := &Source{
			Module:      child,
			StatePath:   opt.SourceStatePath,
			HCLFilePath: tfFilePath,
		}
		migratedResource, err := runner.migrateResource(ctx, src)
		if err != nil {
			return fmt.Errorf("migrate a module %s: %w", child.Address, err)
		}
		if runner.Outputter != nil {
			results = append(results, Result{
				Source:           src,
				MigratedResource: migratedResource,
			})
		}
	}

	if runner.Outputter != nil {
		if err := runner.Outputter.Output(results); err != nil {
			return fmt.Errorf("output the result: %w", err)
		}
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

	if err := runner.Migrator.Migrate(ctx, source, migratedResource); err != nil {
		return nil, fmt.Errorf("migrate a resource: %w", err)
	}
	return migratedResource, nil
}
