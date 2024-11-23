package tfmigrator

import (
	"context"
)

// QuickRunBatch provides CLI interface to run tfmigrator quickly.
// `flag` package is used.
//
//	-help
//	-dry-run
//	-log-level
//	-state - source state file path
//	args - Terraform Configuration file paths
//
// Compared with QuickRun, QuickRunBatch is a low level API.
func QuickRunBatch(ctx context.Context, planner BatchPlanner) error {
	return quickRun(ctx, planner, nil)
}
