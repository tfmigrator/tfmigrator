package tfmigrator

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	tflog "github.com/suzuki-shunsuke/tfmigrator-sdk/tfmigrator/log"
)

// QuickRun provides CLI interface to run tfmigrator quickly.
func QuickRun(ctx context.Context, planner Planner) {
	if err := quickRun(ctx, planner); err != nil {
		log.Fatal(err)
	}
}

func quickRun(ctx context.Context, planner Planner) error {
	logger := &tflog.SimpleLogger{}

	var dryRun bool
	var help bool
	var logLevel string
	flag.BoolVar(&dryRun, "dry-run", false, "dry run")
	flag.BoolVar(&help, "help", false, "show help message")
	flag.StringVar(&logLevel, "log-level", "info", "log level")
	flag.Parse()
	args := flag.Args()

	if help || (len(args) != 0 && args[0] == "help") {
		fmt.Fprint(os.Stderr, `tfmigrator - Migrate Terraform Configuration and State

Usage
  tfmigrator help
  tfmigrator [-help] [-dry-run] [-log-level debug] [Terraform Configuration file path ...]

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

	return runner.Run(ctx, &RunOpt{
		SourceTFFilePaths: args,
	})
}
