package hcledit

import (
	"fmt"
	"io"

	"github.com/minamijoyo/hcledit/editor"
)

type MoveBlockOpt struct {
	From     string
	To       string
	FilePath string
	Stdin    io.Reader
	Stdout   io.Writer
	Update   bool
}

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
