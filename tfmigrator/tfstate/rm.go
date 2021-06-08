package tfstate

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
)

// RemoveOpt is an option of MoveState function.
type RemoveOpt struct {
	StatePath string
}

// Remove runs `terraform state rm`.
func (updater *Updater) Remove(ctx context.Context, address string, opt *RemoveOpt) error {
	cmd := []string{"terraform", "state", "rm"}
	if opt.StatePath != "" {
		cmd = append(cmd, "-state", opt.StatePath)
	}
	cmd = append(cmd, address)
	if updater.DryRun {
		updater.logInfo("[DRYRUN] + " + strings.Join(cmd, " "))
		return nil
	}

	updater.logInfo("+ " + strings.Join(cmd, " "))
	opts := []tfexec.StateRmCmdOption{}
	if opt.StatePath != "" {
		opts = append(opts, tfexec.State(opt.StatePath)) //nolint:staticcheck
	}
	if err := updater.Terraform.StateRm(ctx, address, opts...); err != nil {
		return fmt.Errorf("terraform state rm: %w", err)
	}
	return nil
}
