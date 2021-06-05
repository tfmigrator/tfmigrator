package tfstate

import (
	"context"
)

// MoveOpt is an option of MoveState function.
type MoveOpt struct {
	StatePath string
	StateOut  string
}

func moveStateArgs(sourceAddress, newAddress string, opt *MoveOpt) []string {
	args := []string{"state", "mv"}
	if opt.StatePath != "" {
		args = append(args, "-state", opt.StatePath)
	}
	if opt.StateOut != "" {
		args = append(args, "-state-out", opt.StateOut)
	}

	return append(args, sourceAddress, newAddress)
}

// Move runs `terraform state mv`.
func (updater *Updater) Move(ctx context.Context, sourceAddress, newAddress string, opt *MoveOpt) error {
	return updater.update(ctx, moveStateArgs(sourceAddress, newAddress, opt)...)
}
