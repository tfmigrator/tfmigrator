package tfmigrator

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Songmu/timeout"
)

// MoveStateOpt is an option of MoveState function.
type MoveStateOpt struct {
	StatePath     string
	StateOut      string
	SourceAddress string
	DestAddress   string
}

func moveStateArgs(opt *MoveStateOpt) []string {
	args := []string{"state", "mv"}
	if opt.StatePath != "" {
		args = append(args, "-state", opt.StatePath)
	}
	if opt.StateOut != "" {
		args = append(args, "-state-out", opt.StateOut)
	}

	return append(args, opt.SourceAddress, opt.DestAddress)
}

// MoveState runs `terraform state mv`.
func (runner *Runner) MoveState(ctx context.Context, opt *MoveStateOpt) error {
	args := moveStateArgs(opt)
	if runner.DryRun {
		if runner.Logger != nil {
			runner.Logger.Info("[DRYRUN] + terraform " + strings.Join(args, " "))
		}
		return nil
	}
	if runner.Logger != nil {
		runner.Logger.Info("+ terraform " + strings.Join(args, " "))
	}
	cmd := exec.Command("terraform", args...)
	cmd.Stdout = runner.Stdout
	cmd.Stderr = runner.Stderr
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
