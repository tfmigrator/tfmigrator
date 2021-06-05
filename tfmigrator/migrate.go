package tfmigrator

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
)

// MigrateOpt is an option of Migrate function.
type MigrateOpt struct {
	Stdin      io.Reader
	Stdout     io.Writer
	Stderr     io.Writer
	TFFilePath string
	DryRun     bool
}

// Migrate migrates Terraform Configuration and State with `terraform state mv` and `hcledit`.
func Migrate(ctx context.Context, migratedResource *MigratedResource, opt *MigrateOpt) error {
	// terraform state mv
	if err := MoveState(ctx, &MoveStateOpt{
		StateOut: migratedResource.StatePath(),
		Path:     migratedResource.SourceAddress,
		NewPath:  migratedResource.DestAddress,
		Stderr:   opt.Stderr,
		DryRun:   opt.DryRun,
	}); err != nil {
		return err
	}

	// write tf
	return migrateTF(migratedResource, opt)
}

func migrateTF(migratedResource *MigratedResource, opt *MigrateOpt) error { //nolint:funlen
	if migratedResource.AddressChanged() { //nolint:nestif
		if migratedResource.FileChanged() {
			buf := &bytes.Buffer{}
			if err := getBlock(&getBlockOpt{
				Address: "resource." + migratedResource.SourceAddress,
				File:    opt.TFFilePath,
				Stdout:  buf,
				Stderr:  opt.Stderr,
			}); err != nil {
				return err
			}

			tfPath := migratedResource.TFPath()
			tfFile, err := os.OpenFile(tfPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
			if err != nil {
				return fmt.Errorf("open a file which will write Terraform configuration %s: %w", tfPath, err)
			}
			defer tfFile.Close()

			if err := moveBlock(&moveBlockOpt{
				From:   "resource." + migratedResource.SourceAddress,
				To:     "resource." + migratedResource.DestAddress,
				File:   "-",
				Stdin:  buf,
				Stdout: tfFile,
				Stderr: opt.Stderr,
			}); err != nil {
				return err
			}
			return rmBlock(&rmBlockOpt{
				Address: "resource." + migratedResource.SourceAddress,
				File:    opt.TFFilePath,
				Stdout:  opt.Stdout,
				Stderr:  opt.Stderr,
			})
		}
		return moveBlock(&moveBlockOpt{
			From:   "resource." + migratedResource.SourceAddress,
			To:     "resource." + migratedResource.DestAddress,
			Update: true,
			File:   migratedResource.TFPath(),
			Stdout: opt.Stdout,
			Stderr: opt.Stderr,
		})
	}

	tfPath := migratedResource.TFPath()
	tfFile, err := os.OpenFile(tfPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open a file which will write Terraform configuration %s: %w", tfPath, err)
	}
	defer tfFile.Close()

	if err := getBlock(&getBlockOpt{
		Address: "resource." + migratedResource.SourceAddress,
		File:    opt.TFFilePath,
		Stdout:  tfFile,
		Stderr:  opt.Stderr,
	}); err != nil {
		return err
	}
	return rmBlock(&rmBlockOpt{
		Address: "resource." + migratedResource.SourceAddress,
		File:    opt.TFFilePath,
		Stdout:  opt.Stdout,
		Stderr:  opt.Stderr,
	})
}
