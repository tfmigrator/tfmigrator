package tfmigrator

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func (ctrl *Controller) handleResource(ctx context.Context, param Param, rsc Resource, hclFilePath string, dryRunResult *DryRunResult) error { //nolint:funlen,cyclop
	tfPath := filepath.Join(migratedResource.StateDirname, migratedResource.TFBasename)
	tfFile, err := os.OpenFile(tfPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open a file which will write Terraform configuration %s: %w", tfPath, err)
	}
	defer tfFile.Close()

	hclFile, err := os.Open(hclFilePath)
	if err != nil {
		return fmt.Errorf("open a Terraform configuration %s: %w", hclFilePath, err)
	}
	defer hclFile.Close()

	buf := bytes.Buffer{}
	if err := ctrl.getHCL(ctx, rsc.Address, migratedResource.DestResourcePath, hclFile, &buf); err != nil {
		return err
	}

	if err := ctrl.stateMv(ctx, filepath.Join(migratedResource.StateDirname, migratedResource.StateBasename), rsc.Address, migratedResource.DestResourcePath, param.SkipState); err != nil {
		return err
	}
	// write hcl
	if _, err := io.Copy(tfFile, &buf); err != nil {
		return fmt.Errorf("write Terraform configuration to a file %s: %w", tfPath, err)
	}

	return nil
}
