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
	"github.com/suzuki-shunsuke/tfmigrator-sdk/tfmigrator/log"
)

type Updater struct {
	Stdout io.Writer
	Stderr io.Writer
	DryRun bool
	Logger log.Logger
}

func (updater *Updater) update(ctx context.Context, args ...string) error {
	if updater.DryRun {
		if updater.Logger != nil {
			updater.Logger.Info("[DRYRUN] + terraform " + strings.Join(args, " "))
		}
		return nil
	}
	if updater.Logger != nil {
		updater.Logger.Info("+ terraform " + strings.Join(args, " "))
	}
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
