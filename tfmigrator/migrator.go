package tfmigrator

// Migrator migrates a Terraform resource.
// Note that Migrator doesn't change Terraform State and Terraform Configuration files.
// Migrator determines the updated resource name, outputted State file path, and outputted Terraform Configuration file path.
// If migrator
type Migrator interface {
	Migrate(rsc *Resource) (*MigratedResource, error)
}

// MigrateFunc migrates a Terraform resource.
type MigrateFunc func(rsc *Resource) (*MigratedResource, error)

type newMigrator struct {
	migrate MigrateFunc
}

// NewMigrator is a helper function to create Migrator from MigrateFunc.
func NewMigrator(fn MigrateFunc) Migrator {
	return &newMigrator{
		migrate: fn,
	}
}

func (migrator *newMigrator) Migrate(rsc *Resource) (*MigratedResource, error) {
	return migrator.migrate(rsc)
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

func (migrator *combinedMigrator) Migrate(rsc *Resource) (*MigratedResource, error) {
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
