package common

// Language constants
const (
	LangKotlin = "kotlin"
	LangJava   = "java"
)

// DI Framework constants
const (
	DIHilt = "hilt"
	DIKoin = "koin"
	DINone = "none"
)

// UI Framework constants
const (
	UICompose = "compose"
	UIXML     = "xml"
	UINone    = "none"
)

// Template names
const (
	TemplateEmpty   = "empty"
	TemplateBasic   = "basic"
	TemplateCompose = "compose"
)

// Android defaults
const (
	DefaultMinSdk          = 26
	DefaultTargetSdk       = 35
	DefaultGradleVersion   = "8.9"
	DefaultKotlinVersion   = "1.9.23"
	DefaultComposeVersion  = "1.6.4"
	DefaultJavaTarget      = "1.8"
)

// Spring Boot defaults
const (
	DefaultJavaVersion         = "17"
	DefaultSpringBootVersion   = "3.2.0"
	DefaultSpringBootJavaTarget = "17"
)

// Resource type constants
const (
	ResourceTypeActivity  = "activity"
	ResourceTypeFragment  = "fragment"
	ResourceTypeViewModel = "viewmodel"
)

// Config keys (ProviderConfig)
const (
	ConfigKeyAppName     = "AppName"
	ConfigKeyPackage     = "PackageName"
	ConfigKeyModule      = "Module"
	ConfigKeyLang        = "Lang"
	ConfigKeyTemplate    = "Template"
	ConfigKeyMinSdk      = "MinSdk"
	ConfigKeyTargetSdk   = "TargetSdk"
	ConfigKeyUI          = "UI"
	ConfigKeyDI          = "DI"
	ConfigKeyLayout      = "Layout"
	ConfigKeyKind        = "Kind"
	ConfigKeyName        = "Name"
	ConfigKeyOverride    = "Override"
	ConfigKeyWrapper     = "Wrapper"
	ConfigKeyGitInit     = "Git"
	ConfigKeyGroup       = "Group"
	ConfigKeyArtifact    = "Artifact"
)
