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
	"github.com/sirupsen/logrus"
)

func (ctrl *Controller) tfShow(ctx context.Context, out io.Writer) error {
	logrus.Info("terraform show -json")
	cmd := exec.Command(
		"terraform", "show", "-json")
	cmd.Stdout = out
	cmd.Stderr = ctrl.Stderr
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

func (ctrl *Controller) readStateFromCmd(ctx context.Context, state *State) error {
	buf := bytes.Buffer{}
	if err := ctrl.tfShow(ctx, &buf); err != nil {
		return err
	}
	if err := json.NewDecoder(&buf).Decode(state); err != nil {
		return fmt.Errorf("parse a state as JSON: %w", err)
	}
	return nil
}

func (ctrl *Controller) readState(statePath string, state *State) error {
	stateFile, err := os.Open(statePath)
	if err != nil {
		return fmt.Errorf("open a state file %s: %w", statePath, err)
	}
	defer stateFile.Close()
	if err := json.NewDecoder(stateFile).Decode(state); err != nil {
		return fmt.Errorf("parse a state file as JSON %s: %w", statePath, err)
	}
	return nil
}
