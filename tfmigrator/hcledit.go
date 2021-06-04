package tfmigrator

import (
	"bytes"
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

func (ctrl *Controller) getHCL(
	ctx context.Context, resourcePath, newResourcePath string, hclFile io.Reader, buf io.Writer) error {
	if resourcePath == newResourcePath {
		return ctrl.blockGet(ctx, "resource."+resourcePath, hclFile, buf)
	}
	pp := bytes.Buffer{}
	if err := ctrl.blockGet(ctx, "resource."+resourcePath, hclFile, &pp); err != nil {
		return fmt.Errorf("get a resource from HCL file: %w", err)
	}

	if err := ctrl.blockMv(ctx, "resource."+resourcePath, "resource."+newResourcePath, &pp, buf); err != nil {
		return fmt.Errorf("rename resource: %w", err)
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
