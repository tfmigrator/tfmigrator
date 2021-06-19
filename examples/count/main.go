package main

import (
	"context"
	"fmt"
	"log"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/tfmigrator/tfmigrator/tfmigrator"
)

func main() {
	if err := tfmigrator.QuickRunBatch(context.Background(), tfmigrator.NewBatchPlanner(func(state *tfjson.State, addressFileMap map[string]string) ([]tfmigrator.Result, error) {
		results := []tfmigrator.Result{}
		for _, stateResource := range state.Values.RootModule.Resources {
			results = append(results, tfmigrator.Result{
				Source: &tfmigrator.Source{
					Resource: stateResource,
				},
				MigratedResource: &tfmigrator.MigratedResource{
					Address:          fmt.Sprintf("%s.%s%v", stateResource.Type, stateResource.Name, stateResource.Index),
					SkipHCLMigration: true,
				},
			})
		}
		return results, nil
	})); err != nil {
		log.Fatal(err)
	}
}
