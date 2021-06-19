package tfmigrator

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/tfmigrator/tfmigrator/tfmigrator/hcledit"
	tflog "github.com/tfmigrator/tfmigrator/tfmigrator/log"
	"github.com/tfmigrator/tfmigrator/tfmigrator/tfstate"
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
func QuickRun(ctx context.Context, planner Planner) error { //nolint:funlen
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

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get the current directory: %w", err)
	}
	tfCmdPath, err := exec.LookPath("terraform")
	if err != nil {
		return errors.New("the command `terraform` isn't found: %w")
	}
	tf, err := tfexec.NewTerraform(wd, tfCmdPath)
	if err != nil {
		return fmt.Errorf("initialize Terraform exec: %w", err)
	}

	editor := &hcledit.Client{
		DryRun: dryRun,
		Stderr: os.Stderr,
		Logger: logger,
	}

	runner := &Runner{
		Planner: planner,
		Logger:  logger,
		HCLEdit: editor,
		StateReader: &tfstate.Reader{
			Stderr:    os.Stderr,
			Logger:    logger,
			Terraform: tf,
		},
		Outputter: NewYAMLOutputter(os.Stderr),
		Migrator: &Migrator{
			Stdout:  os.Stdout,
			DryRun:  dryRun,
			HCLEdit: editor,
			StateUpdater: &tfstate.Updater{
				Stdout:    os.Stdout,
				Stderr:    os.Stderr,
				DryRun:    dryRun,
				Logger:    logger,
				Terraform: tf,
			},
		},
		DryRun: dryRun,
	}
	if err := validate.Struct(runner); err != nil {
		return fmt.Errorf("validate Runner: %w", err)
	}

	return runner.Run(ctx, &RunOpt{
		SourceHCLFilePaths: args,
		SourceStatePath:    statePath,
	})
}
