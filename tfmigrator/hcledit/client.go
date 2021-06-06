package hcledit

import (
	"io"

	"github.com/suzuki-shunsuke/tfmigrator/tfmigrator/log"
)

// Client is a client to operate Terraform Configuration files with hcledit.
type Client struct {
	// If DryRun is true, Terraform Configuration files aren't changed.
	// Even if DryRun is true, read operation is done.
	DryRun bool
	// Stderr is an error stream of hcledit's editor.
	Stderr io.Writer
	Logger log.Logger
}

func (client *Client) logDebug(msg string) {
	if client.Logger == nil {
		return
	}
	client.Logger.Debug(msg)
}
