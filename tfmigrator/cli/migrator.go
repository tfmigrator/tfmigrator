package cli

import "github.com/suzuki-shunsuke/tfmigrator-sdk/tfmigrator"

// Migrator migrates a Terraform resource.
// Note that Migrator doesn't change Terraform State and Terraform Configuration files.
// Migrator determines the updated resource name, outputted State file path, and outputted Terraform Configuration file path.
// If migrator
type Migrator interface {
	Migrate(rsc *tfmigrator.Resource) (*tfmigrator.MigratedResource, error)
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

func (migrator *combinedMigrator) Migrate(rsc *tfmigrator.Resource) (*tfmigrator.MigratedResource, error) {
	for _, m := range migrator.migrators {
		migratedResource, err := m.Migrate(rsc)
		if err != nil {
			return nil, err
		}
		if migratedResource == nil {
			continue
		}
		return migratedResource, nil
	}
	return nil, nil
}
