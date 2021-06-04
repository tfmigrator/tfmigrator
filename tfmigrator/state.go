package tfmigrator

type State struct {
	Values Values `json:"values"`
}

type Values struct {
	RootModule RootModule `json:"root_module"`
}

type RootModule struct {
	Resources []Resource `json:"resources"`
}

type Resource struct {
	Address string                 `json:"address"`
	Type    string                 `json:"type"`
	Name    string                 `json:"name"`
	Values  map[string]interface{} `json:"values"`
}
