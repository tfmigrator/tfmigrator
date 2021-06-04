package tfmigrator

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"time"

	"github.com/Songmu/timeout"
	"github.com/sirupsen/logrus"
)

// MoveStateOpt is an option of MoveState function.
type MoveStateOpt struct {
	StateOut string
	Path     string
	NewPath  string
	Stderr   io.Writer
	DryRun   bool
	Logger   *logrus.Entry
}

// MoveState runs `terraform state mv`.
func MoveState(ctx context.Context, opt *MoveStateOpt) error {
	logger := opt.Logger
	if logger == nil {
		logger = logrus.NewEntry(logrus.New())
	}
	if opt.DryRun {
		logger.Info("[DRY RUN] terraform state mv -state-out " + opt.StateOut + " " + opt.Path + " " + opt.NewPath)
		return nil
	}
	logger.Info("terraform state mv -state-out " + opt.StateOut + " " + opt.Path + " " + opt.NewPath)
	cmd := exec.Command("terraform", "state", "mv", "-state-out", opt.StateOut, opt.Path, opt.NewPath)
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
