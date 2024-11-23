package hcledit

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func Test_moveBlock(t *testing.T) { //nolint:paralleltest
	data := []struct {
		title string
		opt   *MoveBlockOpt
		exp   string
	}{
		{
			title: "normal",
			exp: `resource "null_resource" "bar" {}
resource "null_resource" "yoo" {}`,
			opt: &MoveBlockOpt{
				From:     "resource.null_resource.foo",
				To:       "resource.null_resource.yoo",
				FilePath: "-",
				Stdin: strings.NewReader(`resource "null_resource" "bar" {}
resource "null_resource" "foo" {}`),
			},
		},
	}
	client := &Client{
		Stderr: os.Stderr,
	}
	for _, d := range data { //nolint:paralleltest
		t.Run(d.title, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			d.opt.Stdout = stdout
			if err := client.MoveBlock(d.opt); err != nil {
				t.Fatal(err)
			}
			s := stdout.String()
			if d.exp != s {
				t.Fatalf("wanted %s, got %s", d.exp, s)
			}
		})
	}
}
