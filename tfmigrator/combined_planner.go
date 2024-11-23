package tfmigrator

import "fmt"

type combinedPlanner struct {
	plannters []Planner
}

// CombinePlanners creates a migrator from given migrators.
func CombinePlanners(planners ...Planner) Planner {
	return &combinedPlanner{
		plannters: planners,
	}
}

func (cplanner *combinedPlanner) Plan(src *Source) (*MigratedResource, error) {
	for i, p := range cplanner.plannters {
		migratedResource, err := p.Plan(src)
		if err != nil {
			return nil, fmt.Errorf("plan to migrate a resource by combinedPlanner (%d): %w", i, err)
		}
		if migratedResource == nil {
			continue
		}
		return migratedResource, nil
	}
	return nil, nil //nolint:nilnil
}
