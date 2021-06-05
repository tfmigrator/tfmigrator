package main

import (
	"context"
	"flag"
	"log"

	"github.com/suzuki-shunsuke/tfmigrator-sdk/tfmigrator"
)

func main() {
	if err := core(); err != nil {
		log.Fatal(err)
	}
}

func core() error {
	ctx := context.Background()
	logger := &tfmigrator.SimpleLogger{}
	if err := logger.SetLogLevel("debug"); err != nil {
		return err //nolint:wrapcheck
	}

	var dryRun bool
	flag.BoolVar(&dryRun, "dry-run", false, "dry run")
	flag.Parse()

	runner := &tfmigrator.Runner{
		Logger:   logger,
		DryRun:   dryRun,
		Migrator: &migrator{},
	}

	if err := runner.Run(ctx, &tfmigrator.RunOpt{
		SourceTFFilePaths: []string{"main.tf"},
	}); err != nil {
		return err //nolint:wrapcheck
	}
	return nil
}

type migrator struct{}

func (migrator *migrator) Migrate(src *tfmigrator.Source) (*tfmigrator.MigratedResource, error) {
	if src.Address() == "null_resource.foo" {
		return &tfmigrator.MigratedResource{
			Address: "null_resource.bar",
		}, nil
	}
	return nil, nil
}
