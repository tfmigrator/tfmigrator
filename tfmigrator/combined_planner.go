package tfmigrator

import "fmt"

type combinedPlanner struct {
	migrators []Planner
}

// CombinePlanners creates a migrator from given migrators.
func CombinePlanners(migrators ...Planner) Planner {
	return &combinedPlanner{
		migrators: migrators,
	}
}

func (migrator *combinedPlanner) Plan(src *Source) (*MigratedResource, error) {
	for i, m := range migrator.migrators {
		migratedResource, err := m.Plan(src)
		if err != nil {
			return nil, fmt.Errorf("plan to migrate a resource by combinedPlanner (%d): %w", i, err)
		}
		if migratedResource == nil {
			continue
		}
		return migratedResource, nil
	}
	return nil, nil
}
