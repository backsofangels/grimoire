package testing

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/backsofangels/grimoire/internal/providers"
	"github.com/backsofangels/grimoire/internal/validator"
)

// ConfigBuilder provides a fluent interface for building ProviderConfig in tests.
type ConfigBuilder struct {
	cfg providers.ProviderConfig
}

// NewConfig creates a new ConfigBuilder with optional initial key-value pairs.
func NewConfig(kvPairs ...interface{}) *ConfigBuilder {
	cfg := make(providers.ProviderConfig)
	for i := 0; i < len(kvPairs); i += 2 {
		if i+1 < len(kvPairs) {
			k, ok := kvPairs[i].(string)
			if ok {
				cfg[k] = kvPairs[i+1]
			}
		}
	}
	return &ConfigBuilder{cfg: cfg}
}

// Set adds or updates a key-value pair.
func (cb *ConfigBuilder) Set(key string, value interface{}) *ConfigBuilder {
	cb.cfg[key] = value
	return cb
}

// AppName sets the AppName.
func (cb *ConfigBuilder) AppName(name string) *ConfigBuilder {
	cb.cfg["AppName"] = name
	return cb
}

// PackageName sets the PackageName.
func (cb *ConfigBuilder) PackageName(pkg string) *ConfigBuilder {
	cb.cfg["PackageName"] = pkg
	return cb
}

// Lang sets the Lang.
func (cb *ConfigBuilder) Lang(lang string) *ConfigBuilder {
	cb.cfg["Lang"] = lang
	return cb
}

// Module sets the Module.
func (cb *ConfigBuilder) Module(module string) *ConfigBuilder {
	cb.cfg["Module"] = module
	return cb
}

// Kind sets the Kind (for Add operations).
func (cb *ConfigBuilder) Kind(kind string) *ConfigBuilder {
	cb.cfg["Kind"] = kind
	return cb
}

// Name sets the Name (for Add operations).
func (cb *ConfigBuilder) Name(name string) *ConfigBuilder {
	cb.cfg["Name"] = name
	return cb
}

// DI sets the DI (dependency injection) type.
func (cb *ConfigBuilder) DI(di string) *ConfigBuilder {
	cb.cfg["DI"] = di
	return cb
}

// UI sets the UI type.
func (cb *ConfigBuilder) UI(ui string) *ConfigBuilder {
	cb.cfg["UI"] = ui
	return cb
}

// Template sets the Template.
func (cb *ConfigBuilder) Template(template string) *ConfigBuilder {
	cb.cfg["Template"] = template
	return cb
}

// MinSdk sets the MinSdk.
func (cb *ConfigBuilder) MinSdk(sdk int) *ConfigBuilder {
	cb.cfg["MinSdk"] = sdk
	return cb
}

// TargetSdk sets the TargetSdk.
func (cb *ConfigBuilder) TargetSdk(sdk int) *ConfigBuilder {
	cb.cfg["TargetSdk"] = sdk
	return cb
}

// Git sets the Git flag.
func (cb *ConfigBuilder) Git(enabled bool) *ConfigBuilder {
	cb.cfg["Git"] = enabled
	return cb
}

// VSCode sets the Vscode flag.
func (cb *ConfigBuilder) VSCode(enabled bool) *ConfigBuilder {
	cb.cfg["Vscode"] = enabled
	return cb
}

// ViewModel sets the ViewModel flag (for Add operations).
func (cb *ConfigBuilder) ViewModel(enabled bool) *ConfigBuilder {
	cb.cfg["ViewModel"] = enabled
	return cb
}

// Nav sets the Nav flag (for Add operations).
func (cb *ConfigBuilder) Nav(enabled bool) *ConfigBuilder {
	cb.cfg["Nav"] = enabled
	return cb
}

// Override sets the Override flag.
func (cb *ConfigBuilder) Override(enabled bool) *ConfigBuilder {
	cb.cfg["Override"] = enabled
	return cb
}

// Build returns the built ProviderConfig.
func (cb *ConfigBuilder) Build() providers.ProviderConfig {
	return cb.cfg
}

// TestHelper provides common test utilities.
type TestHelper struct {
	t       *testing.T
	tempDir string
}

// NewTestHelper creates a new TestHelper with a temporary directory.
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{
		t:       t,
		tempDir: t.TempDir(),
	}
}

// TempDir returns the temporary directory for the test.
func (th *TestHelper) TempDir() string {
	return th.tempDir
}

// ModulePath returns the path to a module within the test directory.
func (th *TestHelper) ModulePath(moduleName string) string {
	return filepath.Join(th.tempDir, moduleName)
}

// AssertFileExists checks if a file exists at the given path, failing the test if it doesn't.
func (th *TestHelper) AssertFileExists(path string) {
	if _, err := os.Stat(path); err != nil {
		th.t.Fatalf("expected file at %s: %v", path, err)
	}
}

// AssertFileNotExists checks if a file doesn't exist at the given path, failing the test if it does.
func (th *TestHelper) AssertFileNotExists(path string) {
	if _, err := os.Stat(path); err == nil {
		th.t.Fatalf("unexpected file exists at %s", path)
	} else if !os.IsNotExist(err) {
		th.t.Fatalf("error checking file at %s: %v", path, err)
	}
}

// SourcePath returns the expected path for a source file.
func (th *TestHelper) SourcePath(module, pkg, fileName string) string {
	pkgPath := validator.PackageToPath(pkg)
	return filepath.Join(module, "src", "main", "java", pkgPath, fileName)
}

// LayoutPath returns the expected path for a layout XML file.
func (th *TestHelper) LayoutPath(module, layoutName string) string {
	return filepath.Join(module, "src", "main", "res", "layout", layoutName+".xml")
}

// ResourcePath returns the expected path for a resource file.
func (th *TestHelper) ResourcePath(module, resourceType, resourceName string) string {
	return filepath.Join(module, "src", "main", "res", resourceType, resourceName)
}

// ManifestPath returns the expected path for AndroidManifest.xml.
func (th *TestHelper) ManifestPath(module string) string {
	return filepath.Join(module, "src", "main", "AndroidManifest.xml")
}

// BuildGradlePath returns the expected path for build.gradle.
func (th *TestHelper) BuildGradlePath(module string) string {
	return filepath.Join(module, "build.gradle")
}
