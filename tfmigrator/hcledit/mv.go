package hcledit

import (
	"fmt"
	"io"

	"github.com/minamijoyo/hcledit/editor"
)

// MoveBlockOpt is an option of `MoveBlock` function.
type MoveBlockOpt struct {
	// From is a source address.
	From string
	// To is a new address.
	To       string
	FilePath string
	Stdin    io.Reader
	Stdout   io.Writer
	// If `Update` is true, the Terraform Configuration is updated in-place.
	Update bool
}

// MoveBlock moves a block.
func (client *Client) MoveBlock(opt *MoveBlockOpt) error {
	filter := editor.NewBlockRenameFilter(opt.From, opt.To)
	cl := editor.NewClient(&editor.Option{
		InStream:  opt.Stdin,
		OutStream: opt.Stdout,
		ErrStream: client.Stderr,
	})
	if client.DryRun {
		return nil
	}
	if err := cl.Edit(opt.FilePath, opt.Update, filter); err != nil {
		return fmt.Errorf("move a block in %s from %s to %s: %w", opt.FilePath, opt.From, opt.To, err)
	}
	return nil
}
