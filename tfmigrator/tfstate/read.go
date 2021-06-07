package tfstate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"time"

	"github.com/Songmu/timeout"
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
func (reader *Reader) TFShow(ctx context.Context, filePath string, out io.Writer) error {
	args := []string{"show", "-json"}
	if filePath != "" {
		args = append(args, filePath)
	}
	cmd := exec.Command("terraform", args...)
	cmd.Stdout = out
	cmd.Stderr = reader.Stderr
	tioStateMv := timeout.Timeout{
		Cmd:      cmd,
		Duration: 1 * time.Minute,
	}

	msg := "+ terraform show -json"
	if filePath != "" {
		msg += " " + filePath
	}
	reader.logDebug(msg)

	status, err := tioStateMv.RunContext(ctx)
	if err != nil {
		return fmt.Errorf("terraform show -json %s: %w", filePath, err)
	}
	if status.Code != 0 {
		return fmt.Errorf("terraform show -json %s: Exit Code %d", filePath, status.Code)
	}
	return nil
}

// ReadByCmd reads Terraform State by `terraform show -json` command.
func (reader *Reader) ReadByCmd(ctx context.Context, filePath string, state *State) error {
	buf := &bytes.Buffer{}
	if err := reader.TFShow(ctx, filePath, buf); err != nil {
		return err
	}
	return Read(buf, state)
}

// Read reads a file into state.
func Read(file io.Reader, state *State) error {
	if err := json.NewDecoder(file).Decode(state); err != nil {
		return fmt.Errorf("parse a state file as JSON: %w", err)
	}
	return nil
}
