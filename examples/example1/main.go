package main

import (
	"context"

	"github.com/suzuki-shunsuke/tfmigrator/tfmigrator"
)

func main() {
	tfmigrator.QuickRun(context.Background(), tfmigrator.NewPlanner(func(src *tfmigrator.Source) (*tfmigrator.MigratedResource, error) {
		if src.Address() == "null_resource.foo" {
			return &tfmigrator.MigratedResource{
				Address: "null_resource.bar",
			}, nil
		}
		return nil, nil
	}))
}
