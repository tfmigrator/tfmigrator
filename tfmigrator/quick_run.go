package tfmigrator

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	tflog "github.com/tfmigrator/tfmigrator/tfmigrator/log"
)

// QuickRun provides CLI interface to run tfmigrator quickly.
// `flag` package is used.
//   -help
//   -dry-run
//   -log-level
//   -state - source state file path
//   args - Terraform Configuration file paths
// QuickRun is a simple helper function and is designed to implement CLI easily.
// If you want to customize QuickRun, you can use other low level API like `Runner`.
func QuickRun(ctx context.Context, planner Planner) error {
	logger := &tflog.SimpleLogger{}

	var dryRun bool
	var help bool
	var logLevel string
	var statePath string
	flag.BoolVar(&dryRun, "dry-run", false, "dry run")
	flag.BoolVar(&help, "help", false, "show help message")
	flag.StringVar(&logLevel, "log-level", "info", "log level")
	flag.StringVar(&statePath, "state", "", "source State file path")
	flag.Parse()
	args := flag.Args()

	if help || (len(args) != 0 && args[0] == "help") {
		fmt.Fprint(os.Stderr, `tfmigrator - Migrate Terraform Configuration and State

Usage
  tfmigrator help
  tfmigrator [-help] [-dry-run] [-log-level debug] [-state ""] [Terraform Configuration file path ...]

Example

  $ ls *.tf | xargs tfmigrator -dry-run -log-level debug
`)
		return nil
	}

	if len(args) == 0 {
		log.Println("no Terraform Configuration file is passed")
		return nil
	}

	if err := logger.SetLogLevel(logLevel); err != nil {
		return fmt.Errorf("set the log level (%s): %w", logLevel, err)
	}

	runner := &Runner{
		Logger:    logger,
		DryRun:    dryRun,
		Planner:   planner,
		Outputter: NewYAMLOutputter(os.Stderr),
	}
	if err := runner.SetDefault(); err != nil {
		return err
	}
	if err := validate.Struct(runner); err != nil {
		return fmt.Errorf("validate Runner: %w", err)
	}

	return runner.Run(ctx, &RunOpt{
		SourceTFFilePaths: args,
		SourceStatePath:   statePath,
	})
}
