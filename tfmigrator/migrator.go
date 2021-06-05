package tfmigrator

import "fmt"

// Migrator migrates a Terraform resource.
// Note that Migrator doesn't change Terraform State and Terraform Configuration files.
// Migrator determines the updated resource name, outputted State file path, and outputted Terraform Configuration file path.
// If migrator
type Migrator interface {
	Migrate(src *Source) (*MigratedResource, error)
}

// MigrateFunc migrates a Terraform resource.
type MigrateFunc func(src *Source) (*MigratedResource, error)

type newMigrator struct {
	migrate MigrateFunc
}

// NewMigrator is a helper function to create Migrator from MigrateFunc.
func NewMigrator(fn MigrateFunc) Migrator {
	return &newMigrator{
		migrate: fn,
	}
}

func (migrator *newMigrator) Migrate(src *Source) (*MigratedResource, error) {
	return migrator.migrate(src)
}

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
