package tfmigrator

import (
	"fmt"
	"io"

	"github.com/minamijoyo/hcledit/editor"
)

type getBlockOpt struct {
	// e.g. resource.null_resource.foo
	Address string
	// "-", "foo.tf"
	File   string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func getBlock(opt *getBlockOpt) error {
	filter := editor.NewBlockGetFilter(opt.Address)
	client := editor.NewClient(&editor.Option{
		InStream:  opt.Stdin,
		OutStream: opt.Stdout,
		ErrStream: opt.Stderr,
	})
	if err := client.Edit(opt.File, false, filter); err != nil {
		return fmt.Errorf("get a block %s from %s: %w", opt.Address, opt.File, err)
	}
	return nil
}

type moveBlockOpt struct {
	From   string
	To     string
	File   string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func moveBlock(opt *moveBlockOpt) error {
	filter := editor.NewBlockRenameFilter(opt.From, opt.To)
	client := editor.NewClient(&editor.Option{
		InStream:  opt.Stdin,
		OutStream: opt.Stdout,
		ErrStream: opt.Stderr,
	})
	if err := client.Edit(opt.File, false, filter); err != nil {
		return fmt.Errorf("move a block in %s from %s to %s: %w", opt.File, opt.From, opt.To, err)
	}
	return nil
}
