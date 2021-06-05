package tfmigrator

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	tflog "github.com/suzuki-shunsuke/tfmigrator-sdk/tfmigrator/log"
)

// QuickRun provides CLI interface to run tfplanner quickly.
func QuickRun(planner Planner) {
	if err := quickRun(planner); err != nil {
		log.Fatal(err)
	}
}

func quickRun(planner Planner) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
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
		fmt.Fprint(os.Stderr, `tfplanner - Migrate Terraform Configuration and State

Usage
  tfplanner help
  tfplanner [-help] [-dry-run] [-log-level debug] [Terraform Configuration file path ...]

Example

  $ ls *.tf | xargs tfplanner -dry-run -log-level debug
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
		Logger:  logger,
		DryRun:  dryRun,
		Planner: planner,
	}

	return runner.Run(ctx, &RunOpt{
		SourceTFFilePaths: args,
	})
}
