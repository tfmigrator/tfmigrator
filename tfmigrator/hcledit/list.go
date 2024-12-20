package hcledit

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/minamijoyo/hcledit/editor"
)

// ListBlock lists blocks from `filePath` and outputs them into `out`.
func (client *Client) ListBlock(filePath string, out io.Writer) error {
	sink := editor.NewBlockListSink()
	cl := editor.NewClient(&editor.Option{
		OutStream: out,
		ErrStream: client.Stderr,
	})
	client.logDebug("+ hcledit block list -f " + filePath)
	if err := cl.Derive(filePath, sink); err != nil {
		return fmt.Errorf("list blocks in %s: %w", filePath, err)
	}
	return nil
}

// ListBlockMap lists block addresses from `filePath` and returns them as map.
// The key of map is a block address like `resource.null_resource.foo`.
func (client *Client) ListBlockMap(filePath string) (map[string]struct{}, error) {
	m := map[string]struct{}{}
	buf := &bytes.Buffer{}
	if err := client.ListBlock(filePath, buf); err != nil {
		return nil, err
	}
	for _, line := range strings.Split(buf.String(), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "resource.") && !strings.HasPrefix(line, "module.") {
			continue
		}
		if _, ok := m[line]; ok {
			return nil, errors.New("resource address is duplicated: " + line)
		}
		m[line] = struct{}{}
	}
	return m, nil
}

type ListBlockMapsOpt struct {
	FilePaths []string `validate:"required"`
}

// ListBlockMaps returns a map of resource address and Terraform Configuration file path.
func (client *Client) ListBlockMaps(filePaths ...string) (map[string]string, error) {
	m := map[string]string{}
	for _, file := range filePaths {
		addresses, err := client.ListBlockMap(file)
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
