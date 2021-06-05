package tfstate

import (
	"context"
)

func removeStateArgs(address string, opt *RemoveOpt) []string {
	args := []string{"state", "rm"}
	if opt.StatePath != "" {
		args = append(args, "-state", opt.StatePath)
	}

	return append(args, address)
}

// RemoveOpt is an option of MoveState function.
type RemoveOpt struct {
	StatePath string
}

// RemoveState runs `terraform state rm`.
func (updater *Updater) Remove(ctx context.Context, address string, opt *RemoveOpt) error {
	return updater.update(ctx, removeStateArgs(address, opt)...)
}
