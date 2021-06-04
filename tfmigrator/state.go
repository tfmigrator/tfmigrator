package tfmigrator

// State is a Terraform State.
type State struct {
	Values Values `json:"values"`
}

// Values is values of Terraform State.
type Values struct {
	RootModule RootModule `json:"root_module"`
}

// RootModule is a root module of Terraform State.
type RootModule struct {
	Resources []Resource `json:"resources"`
}

// Resource is a resource of Terraform State.
type Resource struct {
	Address string                 `json:"address"`
	Type    string                 `json:"type"`
	Name    string                 `json:"name"`
	Values  map[string]interface{} `json:"values"`
}
