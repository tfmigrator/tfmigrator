package hcledit

import (
	"fmt"
	"io"

	"github.com/minamijoyo/hcledit/editor"
)

func (client *Client) GetBlock(filePath, address string, out io.Writer) error {
	filter := editor.NewBlockGetFilter(address)
	cl := editor.NewClient(&editor.Option{
		OutStream: out,
		ErrStream: client.Stderr,
	})
	if err := cl.Edit(filePath, false, filter); err != nil {
		return fmt.Errorf("get a block %s from %s: %w", address, filePath, err)
	}
	return nil
}
