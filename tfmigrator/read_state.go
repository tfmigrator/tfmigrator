package tfmigrator

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/Songmu/timeout"
)

// TFShowOpt is an option of TFShow funciton.
type TFShowOpt struct {
	Stdout io.Writer
	Stderr io.Writer
}

// TFShow gets Terraform State by `terraform show -json` command.
func TFShow(ctx context.Context, opt *TFShowOpt) error {
	cmd := exec.Command("terraform", "show", "-json")
	cmd.Stdout = opt.Stdout
	cmd.Stderr = opt.Stderr
	tioStateMv := timeout.Timeout{
		Cmd:      cmd,
		Duration: 1 * time.Minute,
	}
	status, err := tioStateMv.RunContext(ctx)
	if err != nil {
		return fmt.Errorf("it failed to run a command: %w", err)
	}
	if status.Code != 0 {
		return errors.New("exit code != 0: " + strconv.Itoa(status.Code))
	}
	return nil
}

// ReadStateByCmdOpt is an option of ReadStateByCmd function.
type ReadStateByCmdOpt struct {
	Stderr io.Writer
}

// ReadStateByCmd reads Terraform State by `terraform show -json` command.
func ReadStateByCmd(ctx context.Context, opt *ReadStateByCmdOpt, state *State) error {
	buf := &bytes.Buffer{}
	if err := TFShow(ctx, &TFShowOpt{
		Stdout: buf,
		Stderr: opt.Stderr,
	}); err != nil {
		return err
	}
	return ReadState(buf, state)
}

// ReadStateFromFile opens a Terraform State file and reads it into state.
func ReadStateFromFile(statePath string, state *State) error {
	stateFile, err := os.Open(statePath)
	if err != nil {
		return fmt.Errorf("open a state file %s: %w", statePath, err)
	}
	defer stateFile.Close()
	return ReadState(stateFile, state)
}

// ReadState reads a file into state.
func ReadState(file io.Reader, state *State) error {
	if err := json.NewDecoder(file).Decode(state); err != nil {
		return fmt.Errorf("parse a state file as JSON: %w", err)
	}
	return nil
}
