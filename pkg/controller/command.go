package controller

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

func (ctrl *Controller) stateMv(ctx context.Context, stateOut, oldPath, newPath string, skipState bool) error {
	if skipState {
		logrus.Info("[DRY RUN] terraform state mv -state-out " + stateOut + " " + oldPath + " " + newPath)
		return nil
	}
	logrus.Info("terraform state mv -state-out " + stateOut + " " + oldPath + " " + newPath)
	cmd := exec.Command(
		"terraform", "state", "mv", "-state-out", stateOut, oldPath, newPath)
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

func (ctrl *Controller) blockGet(ctx context.Context, resourcePath string, hclFile io.Reader, stdout io.Writer) error {
	logrus.Info("+ hcledit block get " + resourcePath)
	cmd := exec.Command("hcledit", "block", "get", resourcePath)
	cmd.Stdin = hclFile
	cmd.Stderr = ctrl.Stderr
	cmd.Stdout = stdout
	tio := timeout.Timeout{
		Cmd:      cmd,
		Duration: 1 * time.Minute,
	}
	status, err := tio.RunContext(ctx)
	if err != nil {
		return fmt.Errorf("it failed to run a command: %w", err)
	}
	if status.Code != 0 {
		return errors.New("exit code != 0: " + strconv.Itoa(status.Code))
	}
	return nil
}

func (ctrl *Controller) blockMv(ctx context.Context, newPath, oldPath string, stdin io.Reader, stdout io.Writer) error {
	logrus.Info("+ hcledit block mv " + newPath + " " + oldPath)
	cmd := exec.Command("hcledit", "block", "mv", newPath, oldPath)
	cmd.Stderr = ctrl.Stderr
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	tio := timeout.Timeout{
		Cmd:      cmd,
		Duration: 1 * time.Minute,
	}
	status, err := tio.RunContext(ctx)
	if err != nil {
		return fmt.Errorf("it failed to run a command: %w", err)
	}
	if status.Code != 0 {
		return errors.New("exit code != 0: " + strconv.Itoa(status.Code))
	}
	return nil
}
