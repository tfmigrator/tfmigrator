package tfmigrator

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/suzuki-shunsuke/tfmigrator-sdk/tfmigrator/hcledit"
	"github.com/suzuki-shunsuke/tfmigrator-sdk/tfmigrator/tfstate"
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

	if src.TFFilePath == "" {
		return nil
	}

	return runner.MigrateTF(src, migratedResource)
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
		if err := os.MkdirAll(filepath.Dir(statePath), 0755); err != nil {
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

func appendFile(filePath string) (io.WriteCloser, error) {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("create a file or append to the file: %w", err)
	}
	return f, nil
}

// MigrateTF migrate Terraform Configuration file.
func (runner *Runner) MigrateTF(src *Source, migratedResource *MigratedResource) error { //nolint:cyclop
	client := runner.HCLEdit
	if migratedResource.Removed {
		return client.RemoveBlock(src.TFFilePath, "resource."+src.Address()) //nolint:wrapcheck
	}

	tfFilePath := migratedResource.TFFilePath()
	if tfFilePath == "" {
		tfFilePath = src.TFFilePath
	}

	if src.Address() != migratedResource.Address && migratedResource.Address != "" { //nolint:nestif
		if src.TFFilePath != migratedResource.TFFilePath() {
			buf := &bytes.Buffer{}
			if err := client.GetBlock(src.TFFilePath, "resource."+src.Address(), buf); err != nil {
				return err //nolint:wrapcheck
			}

			if err := os.MkdirAll(filepath.Dir(tfFilePath), 0755); err != nil {
				return fmt.Errorf("create parent directories of Terraform Configuration file %s: %w", tfFilePath, err)
			}
			tfFile, err := appendFile(tfFilePath)
			if err != nil {
				return fmt.Errorf("open a file which will write Terraform configuration %s: %w", tfFilePath, err)
			}
			defer tfFile.Close()

			if err := client.MoveBlock(&hcledit.MoveBlockOpt{
				From:     "resource." + src.Address(),
				To:       "resource." + migratedResource.Address,
				FilePath: "-",
				Stdin:    buf,
				Stdout:   tfFile,
			}); err != nil {
				return err //nolint:wrapcheck
			}
			return client.RemoveBlock(src.TFFilePath, "resource."+src.Address()) //nolint:wrapcheck
		}
		return client.MoveBlock(&hcledit.MoveBlockOpt{ //nolint:wrapcheck
			From:     "resource." + src.Address(),
			To:       "resource." + migratedResource.Address,
			Update:   true,
			FilePath: migratedResource.TFFilePath(),
			Stdout:   runner.Stdout,
		})
	}

	if err := os.MkdirAll(filepath.Dir(tfFilePath), 0755); err != nil {
		return fmt.Errorf("create parent directories of Terraform Configuration file %s: %w", tfFilePath, err)
	}
	tfFile, err := appendFile(tfFilePath)
	if err != nil {
		return fmt.Errorf("open a file which will write Terraform configuration %s: %w", tfFilePath, err)
	}
	defer tfFile.Close()

	if err := client.GetBlock(src.TFFilePath, "resource."+src.Address(), tfFile); err != nil {
		return err //nolint:wrapcheck
	}
	return client.RemoveBlock(src.TFFilePath, "resource."+src.Address()) //nolint:wrapcheck
}
