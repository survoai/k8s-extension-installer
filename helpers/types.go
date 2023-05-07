package helpers

type Manifest struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Type        string  `yaml:"type"`
	Version     string  `yaml:"version"`
	Author      string  `yaml:"author"`
	Inputs      []Input `yaml:"inputs"`
}

type Input struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Required    bool   `yaml:"required"`
	Description string `yaml:"description"`
	Default     string `yaml:"default"`
}

type ResourceStatus struct {
	Success  bool
	Message  string
	Resource string
	Kind     string
	Path     string
}

type MultiError struct {
	Errors []error
}
