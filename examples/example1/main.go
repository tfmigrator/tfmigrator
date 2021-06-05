package main

import (
	"github.com/suzuki-shunsuke/tfmigrator-sdk/tfmigrator"
)

func main() {
	tfmigrator.QuickRun(tfmigrator.NewPlanner(func(src *tfmigrator.Source) (*tfmigrator.MigratedResource, error) {
		if src.Address() == "null_resource.foo" {
			return &tfmigrator.MigratedResource{
				Address: "null_resource.bar",
			}, nil
		}
		return nil, nil
	}))
}
