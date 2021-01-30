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
	Resources []Resource `json:"resources"`
}

type Resource map[string]interface{}
