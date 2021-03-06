package hcledit

import (
	"fmt"

	"github.com/minamijoyo/hcledit/editor"
)

// RemoveBlock removes a block `address` from a Terraform Configuration file `filePath`.
func (client *Client) RemoveBlock(filePath, address string) error {
	filter := editor.NewBlockRemoveFilter(address)
	cl := editor.NewClient(&editor.Option{
		ErrStream: client.Stderr,
	})

	if client.DryRun {
		client.logDebug("[DRY RUN] + hcledit block rm -u " + address)
		return nil
	}
	client.logDebug("+ hcledit block rm -u " + address)

	if err := cl.Edit(filePath, true, filter); err != nil {
		return fmt.Errorf("get a block %s from %s: %w", address, filePath, err)
	}
	return nil
}
