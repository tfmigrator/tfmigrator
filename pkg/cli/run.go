package cli

import (
	"fmt"

	"github.com/suzuki-shunsuke/tfmigrator/pkg/controller"
	"github.com/urfave/cli/v2"
)

func (runner *Runner) setCLIArg(c *cli.Context, param controller.Param) (controller.Param, error) { //nolint:unparam
	param.StatePath = c.String("state")
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
