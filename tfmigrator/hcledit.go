package tfmigrator

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

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
	Update bool
}

func moveBlock(opt *moveBlockOpt) error {
	filter := editor.NewBlockRenameFilter(opt.From, opt.To)
	client := editor.NewClient(&editor.Option{
		InStream:  opt.Stdin,
		OutStream: opt.Stdout,
		ErrStream: opt.Stderr,
	})
	if err := client.Edit(opt.File, opt.Update, filter); err != nil {
		return fmt.Errorf("move a block in %s from %s to %s: %w", opt.File, opt.From, opt.To, err)
	}
	return nil
}

type listBlockOpt struct {
	File   string `validate:"required"`
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func listBlock(opt *listBlockOpt) error {
	sink := editor.NewBlockListSink()
	client := editor.NewClient(&editor.Option{
		InStream:  opt.Stdin,
		OutStream: opt.Stdout,
		ErrStream: opt.Stderr,
	})
	if err := client.Derive(opt.File, sink); err != nil {
		return fmt.Errorf("list blocks in %s: %w", opt.File, err)
	}
	return nil
}

type listBlockMapOpt struct {
	File   string `validate:"required"`
	Stdin  io.Reader
	Stderr io.Writer
}

func listBlockMap(opt *listBlockMapOpt) (map[string]struct{}, error) {
	m := map[string]struct{}{}
	buf := &bytes.Buffer{}
	if err := listBlock(&listBlockOpt{
		File:   opt.File,
		Stdin:  opt.Stdin,
		Stdout: buf,
		Stderr: opt.Stderr,
	}); err != nil {
		return nil, err
	}
	for _, line := range strings.Split(buf.String(), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if _, ok := m[line]; ok {
			return nil, errors.New("resource address is duplicated: " + line)
		}
		m[line] = struct{}{}
	}
	return m, nil
}

type listBlockMapsOpt struct {
	Files  []string `validate:"required"`
	Stdin  io.Reader
	Stderr io.Writer
}

// listBlockMaps returns a map of resource address and Terraform Configuration file path.
func listBlockMaps(opt *listBlockMapsOpt) (map[string]string, error) {
	m := map[string]string{}
	for _, file := range opt.Files {
		addresses, err := listBlockMap(&listBlockMapOpt{
			File:   file,
			Stdin:  opt.Stdin,
			Stderr: opt.Stderr,
		})
		if err != nil {
			return nil, err
		}
		for address := range addresses {
			if v, ok := m[address]; ok {
				return nil, fmt.Errorf("resource address (%s) is duplicated: %s and %s", address, v, file)
			}
			m[address] = file
		}
	}
	return m, nil
}

type rmBlockOpt struct {
	// e.g. resource.null_resource.foo
	Address string
	// "-", "foo.tf"
	File   string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func rmBlock(opt *rmBlockOpt) error {
	filter := editor.NewBlockRemoveFilter(opt.Address)
	client := editor.NewClient(&editor.Option{
		InStream:  opt.Stdin,
		OutStream: opt.Stdout,
		ErrStream: opt.Stderr,
	})
	if err := client.Edit(opt.File, true, filter); err != nil {
		return fmt.Errorf("get a block %s from %s: %w", opt.Address, opt.File, err)
	}
	return nil
}
