package tfmigrator

import (
	"bytes"
	"context"
	"fmt"
	"os"

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
	// terraform state mv
	destAddress := migratedResource.Address
	if destAddress == "" {
		destAddress = src.Address()
	}

	if migratedResource.Removed {
		if err := runner.StateUpdater.Remove(ctx, src.Address(), &tfstate.RemoveOpt{
			StatePath: src.StatePath,
		}); err != nil {
			return fmt.Errorf("remove state (%s, %s): %w", src.Address(), src.StatePath, err)
		}
	} else {
		if err := runner.StateUpdater.Move(ctx, src.Address(), destAddress, &tfstate.MoveOpt{
			StatePath: src.StatePath,
			StateOut:  migratedResource.StatePath(),
		}); err != nil {
			return fmt.Errorf("move state: %w", err)
		}
	}

	if src.TFFilePath == "" {
		return nil
	}

	// write tf
	return runner.migrateTF(src, migratedResource)
}

func (runner *Runner) migrateTF(src *Source, migratedResource *MigratedResource) error { //nolint:cyclop
	client := runner.HCLEdit
	if migratedResource.Removed {
		return client.RemoveBlock(src.TFFilePath, "resource."+src.Address()) //nolint:wrapcheck
	}

	tfPath := migratedResource.TFPath()
	if tfPath == "" {
		tfPath = src.TFFilePath
	}

	if src.Address() != migratedResource.Address && migratedResource.Address != "" { //nolint:nestif
		if src.TFFilePath != migratedResource.TFPath() {
			buf := &bytes.Buffer{}
			if err := client.GetBlock(src.TFFilePath, "resource."+src.Address(), buf); err != nil {
				return err //nolint:wrapcheck
			}

			tfFile, err := os.OpenFile(tfPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
			if err != nil {
				return fmt.Errorf("open a file which will write Terraform configuration %s: %w", tfPath, err)
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
			FilePath: migratedResource.TFPath(),
			Stdout:   runner.Stdout,
		})
	}

	tfFile, err := os.OpenFile(tfPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open a file which will write Terraform configuration %s: %w", tfPath, err)
	}
	defer tfFile.Close()

	if err := client.GetBlock(src.TFFilePath, "resource."+src.Address(), tfFile); err != nil {
		return err //nolint:wrapcheck
	}
	return client.RemoveBlock(src.TFFilePath, "resource."+src.Address()) //nolint:wrapcheck
}
