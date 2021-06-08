package main

import (
	"context"
	"log"

	"github.com/tfmigrator/tfmigrator/tfmigrator"
)

func main() {
	if err := tfmigrator.QuickRun(context.Background(), tfmigrator.NewPlanner(func(src *tfmigrator.Source) (*tfmigrator.MigratedResource, error) {
		if src.Address() == "module.foo" {
			return &tfmigrator.MigratedResource{
				Address: "module.bar",
			}, nil
		}
		return nil, nil
	})); err != nil {
		log.Fatal(err)
	}
}
