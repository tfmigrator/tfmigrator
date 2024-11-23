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

type cliArgs struct {
	DryRun    bool
	Help      bool
	LogLevel  string
	StatePath string
}

func showCLIHelp() {
	fmt.Fprint(os.Stderr, `tfmigrator - Migrate Terraform Configuration and State

Usage
  tfmigrator help
  tfmigrator [-help] [-dry-run] [-log-level debug] [-state ""] [Terraform Configuration file path ...]

Example

  $ ls *.tf | xargs tfmigrator -dry-run -log-level debug
`)
}

func parseArgs() *cliArgs {
	cArgs := &cliArgs{}
	flag.BoolVar(&cArgs.DryRun, "dry-run", false, "dry run")
	flag.BoolVar(&cArgs.Help, "help", false, "show help message")
	flag.StringVar(&cArgs.LogLevel, "log-level", "info", "log level")
	flag.StringVar(&cArgs.StatePath, "state", "", "source State file path")
	flag.Parse()
	return cArgs
}

// QuickRun provides CLI interface to run tfmigrator quickly.
// `flag` package is used.
//
//	-help
//	-dry-run
//	-log-level
//	-state - source state file path
//	args - Terraform Configuration file paths
//
// QuickRun is a simple helper function and is designed to implement CLI easily.
// If you want to customize QuickRun, you can use other low level API like `Runner`.
func QuickRun(ctx context.Context, planner Planner) error {
	return quickRun(ctx, nil, planner)
}

func quickRun(ctx context.Context, batchPlanner BatchPlanner, planner Planner) error { //nolint:funlen,cyclop
	logger := &tflog.SimpleLogger{}
	cArgs := parseArgs()
	args := flag.Args()

	if cArgs.Help || (len(args) != 0 && args[0] == "help") {
		showCLIHelp()
		return nil
	}

	if len(args) == 0 {
		log.Println("no Terraform Configuration file is passed")
		return nil
	}

	if err := logger.SetLogLevel(cArgs.LogLevel); err != nil {
		return fmt.Errorf("set the log level (%s): %w", cArgs.LogLevel, err)
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
		DryRun: cArgs.DryRun,
		Stderr: os.Stderr,
		Logger: logger,
	}

	if planner == nil {
		runner := &BatchRunner{
			Planner: batchPlanner,
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
				DryRun:  cArgs.DryRun,
				HCLEdit: editor,
				StateUpdater: &tfstate.Updater{
					Stdout:    os.Stdout,
					Stderr:    os.Stderr,
					DryRun:    cArgs.DryRun,
					Logger:    logger,
					Terraform: tf,
				},
			},
			DryRun: cArgs.DryRun,
		}
		if err := validate.Struct(runner); err != nil {
			return fmt.Errorf("validate Runner: %w", err)
		}

		return runner.Run(ctx, &RunOpt{
			SourceHCLFilePaths: args,
			SourceStatePath:    cArgs.StatePath,
		})
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
			DryRun:  cArgs.DryRun,
			HCLEdit: editor,
			StateUpdater: &tfstate.Updater{
				Stdout:    os.Stdout,
				Stderr:    os.Stderr,
				DryRun:    cArgs.DryRun,
				Logger:    logger,
				Terraform: tf,
			},
		},
		DryRun: cArgs.DryRun,
	}
	if err := validate.Struct(runner); err != nil {
		return fmt.Errorf("validate Runner: %w", err)
	}

	return runner.Run(ctx, &RunOpt{
		SourceHCLFilePaths: args,
		SourceStatePath:    cArgs.StatePath,
	})
}
