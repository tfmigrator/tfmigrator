package tfmigrator_test

import (
	"testing"

	"github.com/tfmigrator/tfmigrator/tfmigrator"
)

func TestValidateResourceName(t *testing.T) {
	t.Parallel()
	data := []struct {
		title string
		name  string
		isErr bool
	}{
		{
			title: "normal",
			name:  "foo",
		},
		{
			title: "invalid",
			name:  "foo bar",
			isErr: true,
		},
	}
	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			t.Parallel()
			if err := tfmigrator.ValidateResourceName(d.name); err != nil {
				if !d.isErr {
					t.Fatal(err)
				}
				return
			}
			if d.isErr {
				t.Fatal("error should be returned")
			}
		})
	}
}
