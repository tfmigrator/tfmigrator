package controller

type Config struct {
	Items []Item
}

type Item struct {
	Rule         string
	StateOut     string `yaml:"state_out"`
	ResourceName string `yaml:"resource_name"`
	TFPath       string `yaml:"tf_path"`
}

type Param struct {
	ConfigFilePath string
	LogLevel       string
	StatePath      string
	Items          []Item
	SkipState      bool
}

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
	Adress string                 `json:"address"`
	Type   string                 `json:"type"`
	Name   string                 `json:"name"`
	Values map[string]interface{} `json:"values"`
}
