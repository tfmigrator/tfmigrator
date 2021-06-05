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

func removeStateArgs(opt *RemoveStateOpt) []string {
	args := []string{"state", "rm"}
	if opt.StatePath != "" {
		args = append(args, "-state", opt.StatePath)
	}

	return append(args, opt.Address)
}

// RemoveStateOpt is an option of MoveState function.
type RemoveStateOpt struct {
	StatePath string
	Address   string
}

func (runner *Runner) updateState(ctx context.Context, args ...string) error {
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
		return fmt.Errorf("terraform %s: %w", strings.Join(args, " "), err)
	}
	if status.Code != 0 {
		return errors.New("exit code isn't zero (" + strconv.Itoa(status.Code) + ")")
	}
	return nil
}

// RemoveState runs `terraform state rm`.
func (runner *Runner) RemoveState(ctx context.Context, opt *RemoveStateOpt) error {
	return runner.updateState(ctx, removeStateArgs(opt)...)
}
