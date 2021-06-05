package tfmigrator

import "fmt"

type combinedMigrator struct {
	migrators []Migrator
}

// CombineMigrators creates a migrator from given migrators.
func CombineMigrators(migrators ...Migrator) Migrator {
	return &combinedMigrator{
		migrators: migrators,
	}
}

func (migrator *combinedMigrator) Migrate(src *Source) (*MigratedResource, error) {
	for i, m := range migrator.migrators {
		migratedResource, err := m.Migrate(src)
		if err != nil {
			return nil, fmt.Errorf("plan to migrate a resource by combinedMigrator (%d): %w", i, err)
		}
		if migratedResource == nil {
			continue
		}
		return migratedResource, nil
	}
	return nil, nil
}
