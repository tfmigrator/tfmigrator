package tfstate

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Songmu/timeout"
	"github.com/suzuki-shunsuke/tfmigrator/tfmigrator/log"
)

// Updater updates Terraform State by `terraform state` commands.
type Updater struct {
	Stdout io.Writer
	Stderr io.Writer
	DryRun bool
	Logger log.Logger
}

func (updater *Updater) logInfo(msg string) {
	if updater.Logger != nil {
		updater.Logger.Info(msg)
	}
}

func (updater *Updater) update(ctx context.Context, args ...string) error {
	if updater.DryRun {
		updater.logInfo("[DRYRUN] + terraform " + strings.Join(args, " "))
		return nil
	}
	updater.logInfo("+ terraform " + strings.Join(args, " "))
	cmd := exec.Command("terraform", args...)
	cmd.Stdout = updater.Stdout
	cmd.Stderr = updater.Stderr
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
