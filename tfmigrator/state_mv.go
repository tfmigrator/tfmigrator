package tfmigrator

import (
	"context"
)

// MoveStateOpt is an option of MoveState function.
type MoveStateOpt struct {
	StatePath     string
	StateOut      string
	SourceAddress string
	DestAddress   string
}

func moveStateArgs(opt *MoveStateOpt) []string {
	args := []string{"state", "mv"}
	if opt.StatePath != "" {
		args = append(args, "-state", opt.StatePath)
	}
	if opt.StateOut != "" {
		args = append(args, "-state-out", opt.StateOut)
	}

	return append(args, opt.SourceAddress, opt.DestAddress)
}

// MoveState runs `terraform state mv`.
func (runner *Runner) MoveState(ctx context.Context, opt *MoveStateOpt) error {
	return runner.updateState(ctx, moveStateArgs(opt)...)
}
