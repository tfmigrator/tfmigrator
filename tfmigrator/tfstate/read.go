package tfstate

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/tfmigrator/tfmigrator/tfmigrator/log"
)

// Reader reads Terraform State.
type Reader struct {
	Stderr io.Writer
	Logger log.Logger
}

func (reader *Reader) logDebug(msg string) {
	if reader.Logger == nil {
		return
	}
	reader.Logger.Debug(msg)
}

// TFShow gets Terraform State by `terraform show -json` command.
func (reader *Reader) TFShow(ctx context.Context, filePath string) (*tfjson.State, error) {
	tfCmdPath, err := exec.LookPath("terraform")
	if err != nil {
		return nil, errors.New("the command `terraform` isn't found: %w")
	}
	tf, err := tfexec.NewTerraform(filepath.Dir(filePath), tfCmdPath)
	if err != nil {
		return nil, fmt.Errorf("initialize Terraform exec: %w", err)
	}

	msg := "+ terraform show -json"
	if filePath != "" {
		msg += " " + filePath
		reader.logDebug(msg)
		state, err := tf.ShowStateFile(ctx, filePath)
		if err != nil {
			return nil, fmt.Errorf("terraform show -json %s: %w", filePath, err)
		}
		return state, nil
	}

	reader.logDebug(msg)
	state, err := tf.Show(ctx)
	if err != nil {
		return nil, fmt.Errorf("terraform show -json: %w", err)
	}
	return state, nil
}
