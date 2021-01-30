package cli

import (
	"context"
	"io"

	"github.com/suzuki-shunsuke/tfmigrator/pkg/constant"
	"github.com/urfave/cli/v2"
)

type Runner struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (runner *Runner) Run(ctx context.Context, args ...string) error {
	app := cli.App{
		Name:    "tfmigrator",
		Usage:   "Migrate Terraform Configuration and State. https://github.com/suzuki-shunsuke/tfmigrator",
		Version: constant.Version,
		Commands: []*cli.Command{
			{
				Name:   "run",
				Usage:  "Migrate Terraform Configuration and State",
				Action: runner.runAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "log-level",
						Usage: "log level",
					},
					&cli.BoolFlag{
						Name:  "skip-state",
						Usage: "skip to terraform state mv",
					},
				},
			},
		},
	}

	return app.RunContext(ctx, args)
}
