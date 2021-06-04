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
)

// MoveStateOpt is an option of MoveState function.
type MoveStateOpt struct {
	StateOut string
	Path     string
	NewPath  string
	Stderr   io.Writer
	DryRun   bool
	Logger   Logger
}

// MoveState runs `terraform state mv`.
func MoveState(ctx context.Context, opt *MoveStateOpt) error {
	if opt.DryRun {
		if opt.Logger != nil {
			opt.Logger.Info("[DRYRUN] + terraform state mv -state-out " + opt.StateOut + " " + opt.Path + " " + opt.NewPath)
		}
		return nil
	}
	if opt.Logger != nil {
		opt.Logger.Info("+ terraform state mv -state-out " + opt.StateOut + " " + opt.Path + " " + opt.NewPath)
	}
	cmd := exec.Command("terraform", "state", "mv", "-state-out", opt.StateOut, opt.Path, opt.NewPath) //nolint:gosec
	cmd.Stderr = opt.Stderr
	tioStateMv := timeout.Timeout{
		Cmd:      cmd,
		Duration: 1 * time.Minute,
	}
	status, err := tioStateMv.RunContext(ctx)
	if err != nil {
		return fmt.Errorf("terraform state mv: %w", err)
	}
	if status.Code != 0 {
		return errors.New("exit code of terraform state mv isn't zero (" + strconv.Itoa(status.Code) + ")")
	}
	return nil
}
