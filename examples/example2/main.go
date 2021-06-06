package main

import (
	"context"

	"github.com/suzuki-shunsuke/tfmigrator/tfmigrator"
)

func main() {
	tfmigrator.QuickRun(context.Background(), tfmigrator.NewPlanner(func(src *tfmigrator.Source) (*tfmigrator.MigratedResource, error) {
		if src.Address() == "null_resource.foo" {
			return &tfmigrator.MigratedResource{
				Dirname:        "foo",
				StateBasename:  "terraform.tfstate",
				TFFileBasename: "main.tf",
			}, nil
		}
		if src.Address() == "null_resource.zoo" {
			return &tfmigrator.MigratedResource{
				Dirname:        "foo",
				StateBasename:  "terraform.tfstate",
				TFFileBasename: "main.tf",
			}, nil
		}
		if src.Address() == "null_resource.bar" {
			return &tfmigrator.MigratedResource{
				Removed: true,
			}, nil
		}
		return nil, nil
	}))
}