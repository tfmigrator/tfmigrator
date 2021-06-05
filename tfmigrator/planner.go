package tfmigrator

// Planner plans how a Terraform resource is migrated.
// Note that Planner itself doesn't change Terraform State and Terraform Configuration files.
// Planner determines the updated resource name, outputted State file path, and outputted Terraform Configuration file path.
// tfmigrator migrates according to the plan of Planner.
type Planner interface {
	Plan(src *Source) (*MigratedResource, error)
}

// PlanFunc plans how a Terraform resource is migrated.
type PlanFunc func(src *Source) (*MigratedResource, error)

type newPlanner struct {
	plan PlanFunc
}

// NewPlanner is a helper function to create Planner from PlanFunc.
func NewPlanner(fn PlanFunc) Planner {
	return &newPlanner{
		plan: fn,
	}
}

func (planner *newPlanner) Plan(src *Source) (*MigratedResource, error) {
	return planner.plan(src)
}
