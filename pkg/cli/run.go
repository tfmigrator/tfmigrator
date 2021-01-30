package cli

import (
	"errors"
	"fmt"

	"github.com/suzuki-shunsuke/tfmigrator/pkg/controller"
	"github.com/urfave/cli/v2"
)

func (runner *Runner) setCLIArg(c *cli.Context, param controller.Param) (controller.Param, error) {
	args := c.Args()
	if args.Len() != 1 {
		return controller.Param{}, errors.New(`one arguments are required.
Usage: tfmigrator run <file path to Terraform State>`)
	}
	param.StatePath = args.First()
	if logLevel := c.String("log-level"); logLevel != "" {
		param.LogLevel = logLevel
	}
	param.SkipState = c.Bool("skip-state")
	param.ConfigFilePath = c.String("config")
	if param.ConfigFilePath == "" {
		param.ConfigFilePath = "tfmigrator.yaml"
	}
	return param, nil
}

func (runner *Runner) runAction(c *cli.Context) error {
	param, err := runner.setCLIArg(c, controller.Param{})
	if err != nil {
		return fmt.Errorf("parse the command line arguments: %w", err)
	}

	ctrl, param, err := controller.New(c.Context, param)
	if err != nil {
		return fmt.Errorf("initialize a controller: %w", err)
	}

	return ctrl.Run(c.Context, param)
}
