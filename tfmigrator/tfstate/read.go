package tfstate

import (
	"context"
	"fmt"
	"io"

	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/tfmigrator/tfmigrator/tfmigrator/log"
)

// Reader reads Terraform State.
type Reader struct {
	Stderr    io.Writer
	Logger    log.Logger
	Terraform *tfexec.Terraform
}

func (reader *Reader) logDebug(msg string) {
	if reader.Logger == nil {
		return
	}
	reader.Logger.Debug(msg)
}

// Read gets Terraform State by `terraform show -json` command.
func (reader *Reader) Read(ctx context.Context, filePath string) (*tfjson.State, error) {
	reader.Terraform.SetStderr(reader.Stderr)

	msg := "+ terraform show -json"
	if filePath != "" {
		msg += " " + filePath
		reader.logDebug(msg)
		state, err := reader.Terraform.ShowStateFile(ctx, filePath)
		if err != nil {
			return nil, fmt.Errorf("terraform show -json %s: %w", filePath, err)
		}
		return state, nil
	}

	reader.logDebug(msg)
	state, err := reader.Terraform.Show(ctx)
	if err != nil {
		return nil, fmt.Errorf("terraform show -json: %w", err)
	}
	return state, nil
}
