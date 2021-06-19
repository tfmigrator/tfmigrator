package tfmigrator

import (
	tfjson "github.com/hashicorp/terraform-json"
)

// BatchPlanner plans how Terraform resources are migrated.
// Note that BatchPlanner itself doesn't change Terraform State and Terraform Configuration files.
// BatchPlanner determines the updated resource name, outputted State file path, and outputted Terraform Configuration file path.
// tfmigrator migrates according to the plan of BatchPlanner.
type BatchPlanner interface {
	Plan(state *tfjson.State, addressFileMap map[string]string) ([]Result, error)
}

// BatchPlanFunc plans how Terraform resource are migrated.
type BatchPlanFunc func(state *tfjson.State, addressFileMap map[string]string) ([]Result, error)

type newBatchPlanner struct {
	plan BatchPlanFunc
}

// NewBatchPlanner is a helper function to create BatchPlanner from PlanFunc.
func NewBatchPlanner(fn BatchPlanFunc) BatchPlanner {
	return &newBatchPlanner{
		plan: fn,
	}
}

func (planner *newBatchPlanner) Plan(state *tfjson.State, addressFileMap map[string]string) ([]Result, error) {
	return planner.plan(state, addressFileMap)
}
