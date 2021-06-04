package tfmigrator

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
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
}

// MoveState runs `terraform state mv`.
func MoveState(ctx context.Context, opt *MoveStateOpt) error {
	if opt.DryRun {
		log.Println("[DRY RUN] terraform state mv -state-out " + opt.StateOut + " " + opt.Path + " " + opt.NewPath)
		return nil
	}
	log.Println("terraform state mv -state-out " + opt.StateOut + " " + opt.Path + " " + opt.NewPath)
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
