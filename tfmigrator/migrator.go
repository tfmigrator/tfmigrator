package tfmigrator

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
