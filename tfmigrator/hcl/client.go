package hcl

import (
	"io"

	"github.com/minamijoyo/hcledit/editor"
)

type GetBlockOpt struct {
	Address string
	File    string
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
}

func GetBlock(opt *GetBlockOpt) error {
	filter := editor.NewBlockGetFilter(opt.Address)
	client := editor.NewClient(&editor.Option{
		InStream:  opt.Stdin,
		OutStream: opt.Stdout,
		ErrStream: opt.Stderr,
	})
	return client.Edit(opt.File, false, filter)
}

type MoveBlockOpt struct {
	From   string
	To     string
	File   string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func MoveBlock(opt *MoveBlockOpt) error {
	filter := editor.NewBlockRenameFilter(opt.From, opt.To)
	client := editor.NewClient(&editor.Option{
		InStream:  opt.Stdin,
		OutStream: opt.Stdout,
		ErrStream: opt.Stderr,
	})
	return client.Edit(opt.File, false, filter)
}
