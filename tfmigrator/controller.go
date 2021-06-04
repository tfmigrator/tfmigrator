package controller

import (
	"context"
	"io"
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/go-template-unmarshaler/text"
)

type Controller struct { //nolint:maligned
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func New(ctx context.Context, param Param) (Controller, Param, error) {
	if param.LogLevel != "" {
		lvl, err := logrus.ParseLevel(param.LogLevel)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"log_level": param.LogLevel,
			}).WithError(err).Error("the log level is invalid")
		}
		logrus.SetLevel(lvl)
	}

	text.SetTemplateFunc(func(s string) (*template.Template, error) {
		return template.New("_").Funcs(sprig.TxtFuncMap()).Parse(s) //nolint:wrapcheck
	})

	return Controller{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}, param, nil
}
