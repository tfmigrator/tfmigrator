package tfstate

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
)

// MoveOpt is an option of MoveState function.
type MoveOpt struct {
	StatePath string
	StateOut  string
}

// Move runs `terraform state mv`.
func (updater *Updater) Move(ctx context.Context, sourceAddress, newAddress string, opt *MoveOpt) error {
	cmd := []string{"terraform", "state", "mv"}
	if opt.StatePath != "" {
		cmd = append(cmd, "-state", opt.StatePath)
	}
	if opt.StateOut != "" {
		cmd = append(cmd, "-state-out", opt.StateOut)
	}
	cmd = append(cmd, sourceAddress, newAddress)
	if updater.DryRun {
		updater.logInfo("[DRYRUN] + " + strings.Join(cmd, " "))
		return nil
	}

	opts := []tfexec.StateMvCmdOption{}
	if opt.StatePath != "" {
		opts = append(opts, tfexec.State(opt.StatePath)) //nolint:staticcheck
	}
	if opt.StateOut != "" {
		opts = append(opts, tfexec.StateOut(opt.StateOut))
	}

	updater.logInfo("+ " + strings.Join(cmd, " "))
	if err := updater.Terraform.StateMv(ctx, sourceAddress, newAddress, opts...); err != nil {
		return fmt.Errorf("terraform state mv: %w", err)
	}
	return nil
}
