package tfstate

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

// Reader reads Terraform State.
type Reader struct {
	Stderr io.Writer
}

// TFShow gets Terraform State by `terraform show -json` command.
func (reader *Reader) TFShow(ctx context.Context, out io.Writer) error {
	cmd := exec.Command("terraform", "show", "-json")
	cmd.Stdout = out
	cmd.Stderr = reader.Stderr
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

// ReadByCmd reads Terraform State by `terraform show -json` command.
func (reader *Reader) ReadByCmd(ctx context.Context, state *State) error {
	buf := &bytes.Buffer{}
	if err := reader.TFShow(ctx, buf); err != nil {
		return err
	}
	return Read(buf, state)
}

// ReadFromFile opens a Terraform State file and reads it into state.
func ReadFromFile(statePath string, state *State) error {
	stateFile, err := os.Open(statePath)
	if err != nil {
		return fmt.Errorf("open a state file %s: %w", statePath, err)
	}
	defer stateFile.Close()
	return Read(stateFile, state)
}

// Read reads a file into state.
func Read(file io.Reader, state *State) error {
	if err := json.NewDecoder(file).Decode(state); err != nil {
		return fmt.Errorf("parse a state file as JSON: %w", err)
	}
	return nil
}
