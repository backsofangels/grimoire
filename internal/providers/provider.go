package providers

type ProviderConfig map[string]any

type ProviderFlag struct {
	Name    string
	Short   string
	Usage   string
	Default any
}

type Check struct {
	Name string
	Run  func() error
	Fix  func() error // nil if auto-fix not supported
}

type Provider interface {
	Name() string
	Description() string
	Flags() []ProviderFlag
	Prompt() (ProviderConfig, error)
	Validate(cfg ProviderConfig) error
	Generate(cfg ProviderConfig) error
	// Add creates a single resource inside an existing project (activity/fragment/viewmodel)
	Add(cfg ProviderConfig) error
	DoctorChecks() []Check
}
