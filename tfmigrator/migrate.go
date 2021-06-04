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
	Stderr     io.Writer
	TFFilePath string
	DryRun     bool
}

// Migrate migrates Terraform Configuration and State with `terraform state mv` and `hcledit`.
func Migrate(ctx context.Context, migratedResource *MigratedResource, opt *MigrateOpt) error {
	// terraform state mv
	if err := MoveState(ctx, &MoveStateOpt{
		StateOut: migratedResource.StatePath(),
		Path:     migratedResource.SourceResourcePath,
		NewPath:  migratedResource.DestResourcePath,
		Stderr:   opt.Stderr,
		DryRun:   opt.DryRun,
	}); err != nil {
		return err
	}

	// write tf
	tfPath := migratedResource.TFPath()
	tfFile, err := os.OpenFile(tfPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open a file which will write Terraform configuration %s: %w", tfPath, err)
	}
	defer tfFile.Close()

	if migratedResource.PathChanged() {
		buf := &bytes.Buffer{}
		if err := getBlock(&getBlockOpt{
			Address: "resource." + migratedResource.SourceResourcePath,
			File:    opt.TFFilePath,
			Stdin:   opt.Stdin,
			Stdout:  buf,
			Stderr:  opt.Stderr,
		}); err != nil {
			return err
		}
		return moveBlock(&moveBlockOpt{
			From:   "resource." + migratedResource.SourceResourcePath,
			To:     "resource." + migratedResource.DestResourcePath,
			File:   "-",
			Stdin:  buf,
			Stdout: tfFile,
			Stderr: opt.Stderr,
		})
	}
	return getBlock(&getBlockOpt{
		Address: "resource." + migratedResource.SourceResourcePath,
		File:    opt.TFFilePath,
		Stdin:   opt.Stdin,
		Stdout:  tfFile,
		Stderr:  opt.Stderr,
	})
}
