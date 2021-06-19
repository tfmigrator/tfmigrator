package tfmigrator

import (
	"context"
)

// QuickRunBatch provides CLI interface to run tfmigrator quickly.
// `flag` package is used.
//   -help
//   -dry-run
//   -log-level
//   -state - source state file path
//   args - Terraform Configuration file paths
// QuickRunBatch is a simple helper function and is designed to implement CLI easily.
// If you want to customize QuickRunBatch, you can use other low level API like `Runner`.
func QuickRunBatch(ctx context.Context, planner BatchPlanner) error {
	return quickRun(ctx, planner, nil)
}
