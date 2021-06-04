package tfmigrator

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func Test_getBlock(t *testing.T) { //nolint:paralleltest
	data := []struct {
		title string
		opt   *getBlockOpt
		exp   string
	}{
		{
			title: "normal",
			exp:   `resource "null_resource" "foo" {}`,
			opt: &getBlockOpt{
				Address: "resource.null_resource.foo",
				File:    "-",
				Stdin: strings.NewReader(`resource "null_resource" "bar" {}
resource "null_resource" "foo" {}`),
				Stderr: os.Stderr,
			},
		},
	}
	for _, d := range data { //nolint:paralleltest
		d := d
		t.Run(d.title, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			d.opt.Stdout = stdout
			if err := getBlock(d.opt); err != nil {
				t.Fatal(err)
			}
			s := stdout.String()
			if d.exp != s {
				t.Fatalf("wanted %s, got %s", d.exp, s)
			}
		})
	}
}

func Test_moveBlock(t *testing.T) { //nolint:paralleltest
	data := []struct {
		title string
		opt   *moveBlockOpt
		exp   string
	}{
		{
			title: "normal",
			exp: `resource "null_resource" "bar" {}
resource "null_resource" "yoo" {}`,
			opt: &moveBlockOpt{
				From: "resource.null_resource.foo",
				To:   "resource.null_resource.yoo",
				File: "-",
				Stdin: strings.NewReader(`resource "null_resource" "bar" {}
resource "null_resource" "foo" {}`),
				Stderr: os.Stderr,
			},
		},
	}
	for _, d := range data { //nolint:paralleltest
		d := d
		t.Run(d.title, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			d.opt.Stdout = stdout
			if err := moveBlock(d.opt); err != nil {
				t.Fatal(err)
			}
			s := stdout.String()
			if d.exp != s {
				t.Fatalf("wanted %s, got %s", d.exp, s)
			}
		})
	}
}
