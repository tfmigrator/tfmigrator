package hcledit

import (
	"io"
)

type Client struct {
	DryRun bool
	Stderr io.Writer
}
