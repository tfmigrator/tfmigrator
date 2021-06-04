package tfmigrator

import (
	"io"

	"github.com/minamijoyo/hcledit/editor"
)

type getBlockOpt struct {
	Address string
	File    string
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
}

func getBlock(opt *getBlockOpt) error {
	filter := editor.NewBlockGetFilter(opt.Address)
	client := editor.NewClient(&editor.Option{
		InStream:  opt.Stdin,
		OutStream: opt.Stdout,
		ErrStream: opt.Stderr,
	})
	return client.Edit(opt.File, false, filter)
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
	return client.Edit(opt.File, false, filter)
}
