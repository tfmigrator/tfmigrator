package tfmigrator

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/tfmigrator/tfmigrator/tfmigrator/hcledit"
	"github.com/tfmigrator/tfmigrator/tfmigrator/tfstate"
)

// Migrate migrates Terraform Configuration and State with `terraform state mv` and `hcledit`.
func (runner *Runner) Migrate(ctx context.Context, src *Source, migratedResource *MigratedResource) error {
	if migratedResource == nil {
		return nil
	}
	if err := validate.Struct(src); err != nil {
		return fmt.Errorf("validate Source: %w", err)
	}
	if err := validate.Struct(migratedResource); err != nil {
		return fmt.Errorf("validate MigratedResource: %w", err)
	}

	if err := runner.MigrateState(ctx, src, migratedResource); err != nil {
		return err
	}

	if src.HCLFilePath == "" {
		return nil
	}

	return runner.MigrateHCL(src, migratedResource)
}

// MigrateState migrates Terraform State.
func (runner *Runner) MigrateState(ctx context.Context, src *Source, migratedResource *MigratedResource) error {
	// terraform state mv
	newAddress := migratedResource.Address
	if newAddress == "" {
		newAddress = src.Address()
	}

	if migratedResource.Removed {
		if err := runner.StateUpdater.Remove(ctx, src.Address(), &tfstate.RemoveOpt{
			StatePath: src.StatePath,
		}); err != nil {
			return fmt.Errorf("remove state (%s, %s): %w", src.Address(), src.StatePath, err)
		}
		return nil
	}
	statePath := migratedResource.StatePath()
	if statePath != "" {
		if err := runner.mkdirAll(filepath.Dir(statePath)); err != nil {
			return fmt.Errorf("create parent directories of Terraform State %s: %w", statePath, err)
		}
	}
	if err := runner.StateUpdater.Move(ctx, src.Address(), newAddress, &tfstate.MoveOpt{
		StatePath: src.StatePath,
		StateOut:  statePath,
	}); err != nil {
		return fmt.Errorf("move state: %w", err)
	}
	return nil
}

func (runner *Runner) appendFile(filePath string) (io.WriteCloser, error) {
	if runner.DryRun {
		return &nopWriteCloser{}, nil
	}
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644) //nolint:gomnd
	if err != nil {
		return nil, fmt.Errorf("create a file or append to the file: %w", err)
	}
	return f, nil
}

func (runner *Runner) mkdirAll(p string) error {
	if runner.DryRun {
		return nil
	}
	return os.MkdirAll(p, 0o755) //nolint:wrapcheck,gomnd
}

// MigrateHCL migrate Terraform Configuration file.
func (runner *Runner) MigrateHCL(src *Source, migratedResource *MigratedResource) error { //nolint:cyclop,funlen
	client := runner.HCLEdit
	if migratedResource.Removed {
		return client.RemoveBlock(src.HCLFilePath, "resource."+src.Address()) //nolint:wrapcheck
	}

	tfFilePath := migratedResource.HCLFilePath()
	if tfFilePath == "" {
		tfFilePath = src.HCLFilePath
	}

	if src.Address() != migratedResource.Address && migratedResource.Address != "" { //nolint:nestif
		// address is changed and isn't empty
		if src.HCLFilePath == migratedResource.HCLFilePath() || migratedResource.HCLFilePath() == "" {
			// Terraform Configuration file path isn't changed.
			filePath := migratedResource.HCLFilePath()
			if filePath == "" {
				filePath = src.HCLFilePath
			}
			return client.MoveBlock(&hcledit.MoveBlockOpt{ //nolint:wrapcheck
				From:     src.HCLAddress(),
				To:       migratedResource.HCLAddress(),
				Update:   true,
				FilePath: filePath,
				Stdout:   runner.Stdout,
			})
		}
		// Terraform Configuration file path is changed.
		buf := &bytes.Buffer{}
		if err := client.GetBlock(src.HCLFilePath, "resource."+src.Address(), buf); err != nil {
			return err //nolint:wrapcheck
		}

		if err := runner.mkdirAll(filepath.Dir(tfFilePath)); err != nil {
			return fmt.Errorf("create parent directories of Terraform Configuration file %s: %w", tfFilePath, err)
		}
		tfFile, err := runner.appendFile(tfFilePath)
		if err != nil {
			return fmt.Errorf("open a file which will write Terraform configuration %s: %w", tfFilePath, err)
		}
		defer tfFile.Close()

		if err := client.MoveBlock(&hcledit.MoveBlockOpt{
			From:     src.HCLAddress(),
			To:       migratedResource.HCLAddress(),
			FilePath: "-",
			Stdin:    buf,
			Stdout:   tfFile,
		}); err != nil {
			return err //nolint:wrapcheck
		}
		return client.RemoveBlock(src.HCLFilePath, "resource."+src.Address()) //nolint:wrapcheck
	}

	if err := runner.mkdirAll(filepath.Dir(tfFilePath)); err != nil {
		return fmt.Errorf("create parent directories of Terraform Configuration file %s: %w", tfFilePath, err)
	}
	tfFile, err := runner.appendFile(tfFilePath)
	if err != nil {
		return fmt.Errorf("open a file which will write Terraform configuration %s: %w", tfFilePath, err)
	}
	defer tfFile.Close()

	if err := client.GetBlock(src.HCLFilePath, "resource."+src.Address(), tfFile); err != nil {
		return err //nolint:wrapcheck
	}
	return client.RemoveBlock(src.HCLFilePath, "resource."+src.Address()) //nolint:wrapcheck
}
