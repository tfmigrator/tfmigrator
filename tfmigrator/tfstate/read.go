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
func (reader *Reader) TFShow(ctx context.Context, out io.Writer) error {
	cmd := exec.Command("terraform", "show", "-json")
	cmd.Stdout = out
	cmd.Stderr = reader.Stderr
	tioStateMv := timeout.Timeout{
		Cmd:      cmd,
		Duration: 1 * time.Minute,
	}

	reader.logDebug("+ terraform show -json")

	status, err := tioStateMv.RunContext(ctx)
	if err != nil {
		return fmt.Errorf("terraform show -json: %w", err)
	}
	if status.Code != 0 {
		return errors.New("terraform show -json: exit code != 0: " + strconv.Itoa(status.Code))
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
