package hcledit

import (
	"fmt"
	"io"

	"github.com/minamijoyo/hcledit/editor"
)

// GetBlock gets a block `address` from `filePath` and outputs it into `out`.
// address is an address like `resource.null_resource.foo`.
func (client *Client) GetBlock(filePath, address string, out io.Writer) error {
	filter := editor.NewBlockGetFilter(address)
	cl := editor.NewClient(&editor.Option{
		OutStream: out,
		ErrStream: client.Stderr,
	})
	client.logDebug(fmt.Sprintf("+ hcledit block get -f %s %s", filePath, address))
	if err := cl.Edit(filePath, false, filter); err != nil {
		return fmt.Errorf("get a block %s from %s: %w", address, filePath, err)
	}
	return nil
}
