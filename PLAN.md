# Grimoire Refactoring Plan

**Objective:** Improve code clarity, reusability, and readability across the Grimoire codebase by eliminating 300+ lines of duplication and consolidating boilerplate.

**Timeline:** ~10-14 hours (4 phases, can be parallelized where noted)
**Risk Level:** Low (all changes backward compatible, tested incrementally)
**Current Status:** Analysis phase complete; ready for implementation

---

## Phase 1: Quick Wins & Foundation (1-2 hours)

These quick wins establish the foundation for later refactorings and are low-risk, high-impact tasks.

### 1.1 Remove Unused Imports & Format Code
**Goal:** Clean up codebase; establish consistent formatting baseline.

**Files to audit:**
- cmd/*.go
- internal/providers/**/*.go
- internal/validator/*.go
- internal/config/*.go
- internal/logging/*.go

**Tasks:**
1. Run `gofmt -s ./...` across entire repo
2. Run `go mod tidy` to remove unused dependencies
3. Use `go vet ./...` to detect unused variables/imports
4. Manually audit and remove unused imports in all files
5. Verify all tests still pass: `go test ./...`

**Verification:**
```bash
go vet ./...
go fmt ./...
go test ./...
```

**Estimated Effort:** 20 minutes
**Affected Files:** All Go files (no logic changes)

---

### 1.2 Define String & Numeric Constants
**Goal:** Replace magic strings and numbers with named constants.

**Files to create:**
- `internal/common/constants.go` (new)

**Tasks:**
1. Create `internal/common/` directory
2. Create `constants.go` with:
   ```go
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
   )
   
   // Spring Boot defaults
   const (
       DefaultJavaVersion     = "17"
       DefaultSpringBootVersion = "3.2.0"
   )
   ```

3. Find all hardcoded strings and magic numbers:
   - `if lang == "kotlin"` → `if lang == LangKotlin`
   - `if minSdk == 0 { minSdk = 26 }` → `if minSdk == 0 { minSdk = DefaultMinSdk }`
   - Replace in files:
     - `internal/providers/android/provider.go`
     - `internal/providers/android/generator.go`
     - `internal/providers/android/add.go`
     - `internal/providers/springboot/generator.go`
     - `cmd/add.go`

4. Update imports in all affected files
5. Run tests: `go test ./...`

**Verification:**
```bash
go test ./...
grep -r "== \"kotlin\"" internal/ cmd/  # Should be minimal/zero
grep -r "== \"compose\"" internal/ cmd/  # Should be minimal/zero
```

**Estimated Effort:** 45 minutes
**Affected Files:** 5-6 provider & command files

---

### 1.3 Create Normalization Utility (internal/common/normalize.go)
**Goal:** Centralize string normalization logic; replace 20+ scattered instances.

**File to create:**
- `internal/common/normalize.go`

**Tasks:**
1. Create `internal/common/normalize.go`:
   ```go
   package common
   
   import "strings"
   
   // Normalize applies lowercase and trimspace.
   func Normalize(s string) string {
       return strings.ToLower(strings.TrimSpace(s))
   }
   
   // NormalizeLang normalizes a language identifier.
   func NormalizeLang(s string) string {
       return Normalize(s)
   }
   
   // NormalizeUI normalizes a UI type identifier.
   func NormalizeUI(s string) string {
       return Normalize(s)
   }
   
   // NormalizeDI normalizes a DI framework identifier.
   func NormalizeDI(s string) string {
       return Normalize(s)
   }
   ```

2. Replace all instances of `strings.ToLower(strings.TrimSpace(...))` in:
   - `cmd/add.go` (~5 instances)
   - `internal/providers/android/provider.go` (~3 instances)
   - `internal/providers/android/add.go` (~4 instances)

3. Run tests: `go test ./...`

**Verification:**
```bash
go test ./...
grep -r "strings.ToLower(strings.TrimSpace" internal/ cmd/ | wc -l  # Should be 0-1
```

**Estimated Effort:** 30 minutes
**Affected Files:** 3-4 files

---

## Phase 2: Config & Validation Helpers (2-3 hours)

These establish reusable patterns for configuration extraction and validation, reducing boilerplate by ~100 lines.

### 2.1 Create Config Extraction Helpers (internal/common/config_helpers.go)
**Goal:** Replace 100+ lines of duplicated config extraction/fallback logic.

**File to create:**
- `internal/common/config_helpers.go`

**Tasks:**
1. Create `internal/common/config_helpers.go`:
   ```go
   package common
   
   import (
       "strconv"
       "github.com/backsofangels/grimoire/internal/providers"
   )
   
   // GetString retrieves a string config value with fallback to multiple keys.
   // Returns empty string if not found.
   func GetString(cfg providers.ProviderConfig, keys ...string) string {
       for _, k := range keys {
           if v, ok := cfg[k].(string); ok && v != "" {
               return v
           }
       }
       return ""
   }
   
   // GetStringDefault retrieves a string config value with a default fallback.
   func GetStringDefault(cfg providers.ProviderConfig, defaultVal string, keys ...string) string {
       if s := GetString(cfg, keys...); s != "" {
           return s
       }
       return defaultVal
   }
   
   // GetInt retrieves an int config value with type coercion from string.
   func GetInt(cfg providers.ProviderConfig, keys ...string) int {
       for _, k := range keys {
           if v, ok := cfg[k].(int); ok && v > 0 {
               return v
           }
           if v, ok := cfg[k].(string); ok {
               if n, err := strconv.Atoi(v); err == nil && n > 0 {
                   return n
               }
           }
       }
       return 0
   }
   
   // GetIntDefault retrieves an int config value with a default fallback.
   func GetIntDefault(cfg providers.ProviderConfig, defaultVal int, keys ...string) int {
       if n := GetInt(cfg, keys...); n > 0 {
           return n
       }
       return defaultVal
   }
   
   // GetBool retrieves a bool config value.
   func GetBool(cfg providers.ProviderConfig, keys ...string) bool {
       for _, k := range keys {
           if v, ok := cfg[k].(bool); ok {
               return v
           }
           if v, ok := cfg[k].(string); ok && (v == "true" || v == "1") {
               return true
           }
       }
       return false
   }
   ```

2. Update `internal/providers/android/generator.go`:
   - Replace all config extraction boilerplate with calls to helpers
   - Example replacements:
     ```go
     // Before (5-6 lines)
     lang, _ := cfg["Lang"].(string)
     if lang == "" {
         if l2, _ := cfg["lang"].(string); l2 != "" {
             lang = l2
         } else {
             lang = "kotlin"
         }
     }
     
     // After (1 line)
     lang := common.GetStringDefault(cfg, constants.LangKotlin, "Lang", "lang")
     ```

3. Update `internal/providers/springboot/generator.go` similarly (~40 lines reduced)

4. Add unit tests for helpers in `internal/common/config_helpers_test.go`

5. Run tests: `go test ./...`

**Verification:**
```bash
go test ./... 
# Count lines with fallback pattern in generators (should be near 0)
grep -c "if.*ok :=.*cfg\[" internal/providers/android/generator.go
grep -c "if.*ok :=.*cfg\[" internal/providers/springboot/generator.go
```

**Estimated Effort:** 1 hour
**Affected Files:** android/generator.go, springboot/generator.go, new config_helpers.go

---

### 2.2 Consolidate Validation Functions (expand internal/validator/validator.go)
**Goal:** Move scattered validation logic into single validator package.

**Files involved:**
- `cmd/add.go` → move validateUI, validateLang, validateDI
- `internal/providers/android/add.go` → move validateClassName
- Consolidate into `internal/validator/validator.go`

**Tasks:**
1. In `internal/validator/validator.go`, add functions:
   ```go
   // ValidateLanguage checks if language is valid (kotlin|java).
   func ValidateLanguage(lang string) error {
       s := strings.ToLower(strings.TrimSpace(lang))
       if s == "" {
           return nil  // Empty is OK (uses default)
       }
       switch s {
       case "kotlin", "java":
           return nil
       default:
           return fmt.Errorf("invalid language: %s (allowed: kotlin|java)", lang)
       }
   }
   
   // ValidateUI checks if UI type is valid (xml|compose).
   func ValidateUI(ui string) error {
       s := strings.ToLower(strings.TrimSpace(ui))
       if s == "" || s == "none" {
           return nil  // Empty or 'none' is OK (disables UI)
       }
       switch s {
       case "xml", "compose":
           return nil
       default:
           return fmt.Errorf("invalid UI type: %s (allowed: xml|compose)", ui)
       }
   }
   
   // ValidateDI checks if DI framework is valid (none|hilt|koin).
   func ValidateDI(di string) error {
       s := strings.ToLower(strings.TrimSpace(di))
       if s == "" {
           return nil  // Empty is OK (uses default)
       }
       switch s {
       case "none", "hilt", "koin":
           return nil
       default:
           return fmt.Errorf("invalid DI framework: %s (allowed: none|hilt|koin)", di)
       }
   }
   
   // ValidateClassName checks if a class name is valid Kotlin/Java.
   func ValidateClassName(name string) error {
       if name == "" {
           return fmt.Errorf("class name cannot be empty")
       }
       if !isValidIdentifier(name) {
           return fmt.Errorf("invalid class name: %s (must start with letter and contain only alphanumeric/underscore)", name)
       }
       if len(name) > 128 {
           return fmt.Errorf("class name too long (max 128 characters)")
       }
       return nil
   }
   ```

2. Add unit tests in `internal/validator/validator_test.go`

3. Update imports in:
   - `cmd/add.go`: replace local validateUI/validateLang/validateDI with validator.ValidateUI/Language/DI
   - `internal/providers/android/add.go`: import validator; use validator.ValidateClassName

4. Remove local validation functions from cmd/add.go (~30 lines removed)

5. Run tests: `go test ./...`

**Verification:**
```bash
go test ./...
grep -c "func validateUI\|func validateLang\|func validateDI" cmd/add.go  # Should be 0
```

**Estimated Effort:** 45 minutes
**Affected Files:** cmd/add.go, android/add.go, validator/validator.go, validator/validator_test.go (new tests)

---

### 2.3 Create Common Generator Helpers (internal/common/generator.go)
**Goal:** Extract renderTemplate, writeFile, initGit from both providers; eliminate 40+ duplicate lines.

**File to create:**
- `internal/common/generator.go`

**Tasks:**
1. Create `internal/common/generator.go`:
   ```go
   package common
   
   import (
       "fmt"
       "os"
       "os/exec"
       "path/filepath"
       "text/template"
   )
   
   // RenderTemplate renders a Go template string with the given data.
   func RenderTemplate(tmplName, tmplContent string, data any) (string, error) {
       t, err := template.New(tmplName).Parse(tmplContent)
       if err != nil {
           return "", fmt.Errorf("parse template %s: %w", tmplName, err)
       }
       var buf bytes.Buffer
       if err := t.Execute(&buf, data); err != nil {
           return "", fmt.Errorf("execute template %s: %w", tmplName, err)
       }
       return buf.String(), nil
   }
   
   // WriteFile writes content to a file, creating directories as needed.
   func WriteFile(path, content string) error {
       if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
           return fmt.Errorf("mkdir %s: %w", filepath.Dir(path), err)
       }
       if err := os.WriteFile(path, []byte(content), 0644); err != nil {
           return fmt.Errorf("write %s: %w", path, err)
       }
       return nil
   }
   
   // InitGit initializes a git repository in the given directory.
   func InitGit(dir string) error {
       cmd := exec.Command("git", "init")
       cmd.Dir = dir
       if err := cmd.Run(); err != nil {
           return fmt.Errorf("git init: %w", err)
       }
       return nil
   }
   ```

2. Update `internal/providers/android/generator.go`:
   - Remove local renderTemplate, writeFile, initGit functions (~40 lines)
   - Import `"github.com/backsofangels/grimoire/internal/common"`
   - Replace calls: `renderTemplate()` → `common.RenderTemplate()`, etc.

3. Update `internal/providers/springboot/generator.go` similarly

4. Create unit tests in `internal/common/generator_test.go`

5. Run tests: `go test ./...`

**Verification:**
```bash
go test ./...
# Verify duplicate functions removed
grep -c "^func renderTemplate" internal/providers/android/generator.go
grep -c "^func renderTemplate" internal/providers/springboot/generator.go
# Both should be 0
```

**Estimated Effort:** 50 minutes
**Affected Files:** android/generator.go, springboot/generator.go, new common/generator.go

---

## Phase 3: Boilerplate Reduction & Refactoring (3-4 hours)

These refactorings consolidate repeated patterns and break down large functions.

### 3.1 Extract Add Command Handler Logic (cmd/add_helpers.go)
**Goal:** Replace ~70 duplicate lines from three add subcommand handlers with single orchestration function.

**File to create:**
- `cmd/add_helpers.go`

**Tasks:**
1. Create `cmd/add_helpers.go` with:
   ```go
   package cmd
   
   import (
       "fmt"
       "strings"
       "github.com/backsofangels/grimoire/internal/providers"
       "github.com/backsofangels/grimoire/internal/validator"
   )
   
   var ErrNeedsInteractive = fmt.Errorf("interactive mode required")
   
   type AddResourceInput struct {
       Kind       string
       Name       string
       Package    string
       Module     string
       Language   string
       Layout     string
       UI         string
       Override   bool
       DI         string
       ViewModel  bool
       Navigation bool
   }
   
   // extractAddFlags consolidates flag extraction for all add commands.
   func extractAddFlags(cmd *cobra.Command, args []string, kind string) (*AddResourceInput, error) {
       nameFlag, _ := cmd.Flags().GetString("name")
       var name string
       if len(args) == 0 && nameFlag == "" {
           return nil, ErrNeedsInteractive
       }
       if nameFlag != "" {
           name = nameFlag
       } else {
           name = args[0]
       }
       
       // Extract all flags
       pkg, _ := cmd.Flags().GetString("package")
       module, _ := cmd.Flags().GetString("module")
       lang, _ := cmd.Flags().GetString("lang")
       layout, _ := cmd.Flags().GetString("layout")
       ui, _ := cmd.Flags().GetString("ui")
       noUI, _ := cmd.Flags().GetBool("no-ui")
       override, _ := cmd.Flags().GetBool("override")
       di, _ := cmd.Flags().GetString("di")
       
       if noUI {
           ui = "none"
       }
       
       // Validate
       if strings.ToLower(ui) == "none" && layout != "" {
           logging.Info("Ignoring --layout because --ui is 'none'")
           layout = ""
       }
       if strings.ToLower(ui) != "none" {
           if err := validator.ValidateUI(ui); err != nil {
               return nil, fmt.Errorf("invalid --ui: %w", err)
           }
       }
       if err := validator.ValidateLanguage(lang); err != nil {
           return nil, fmt.Errorf("invalid --lang: %w", err)
       }
       if err := validator.ValidateDI(di); err != nil {
           return nil, fmt.Errorf("invalid --di: %w", err)
       }
       
       return &AddResourceInput{
           Kind:      kind,
           Name:      name,
           Package:   pkg,
           Module:    module,
           Language:  lang,
           Layout:    layout,
           UI:        ui,
           Override:  override,
           DI:        di,
       }, nil
   }
   
   // runAddResourceCommand orchestrates the add command flow.
   func runAddResourceCommand(cmd *cobra.Command, provider providers.Provider, kind string, args []string) error {
       input, err := extractAddFlags(cmd, args, kind)
       if err == ErrNeedsInteractive {
           return runAddInteractive(cmd, provider, kind)
       }
       if err != nil {
           return err
       }
       
       cfg := inputToProviderConfig(input)
       if err := provider.Add(cfg); err != nil {
           return fmt.Errorf("add %s failed: %w", kind, err)
       }
       logging.Success(fmt.Sprintf("Added %s %s to %s", kind, input.Name, input.Module))
       return nil
   }
   
   // inputToProviderConfig converts AddResourceInput to ProviderConfig.
   func inputToProviderConfig(input *AddResourceInput) providers.ProviderConfig {
       return providers.ProviderConfig{
           "Kind":        input.Kind,
           "Name":        input.Name,
           "PackageName": input.Package,
           "Module":      input.Module,
           "Lang":        input.Language,
           "Layout":      input.Layout,
           "UI":          input.UI,
           "Override":    input.Override,
           "DI":          input.DI,
       }
   }
   ```

2. Update `cmd/add.go`:
   - Simplify `addActivityCmd`:
     ```go
     var addActivityCmd = &cobra.Command{
         Use:   "activity [name]",
         Short: "Add an Activity to the Android project",
         RunE: func(cmd *cobra.Command, args []string) error {
             provider, _ := providers.Get(cmd.Flag("provider").Value.String())
             return runAddResourceCommand(cmd, provider, "activity", args)
         },
     }
     ```
   - Apply same pattern to `addFragmentCmd` and `addViewModelCmd`
   - Remove all duplicated flag extraction and validation (~70 lines removed per command)

3. Update `cmd/add_prompts.go`:
   - In `runAddInteractive()`, use validators from validator package instead of local functions
   - Add conditional validation: `if ui != "none" { validate }`

4. Run tests: `go test ./...`

**Verification:**
```bash
go test ./...
wc -l cmd/add.go  # Should be reduced from ~400 to ~250
```

**Estimated Effort:** 1 hour
**Affected Files:** cmd/add.go, cmd/add_helpers.go (new), cmd/add_prompts.go

---

### 3.2 Create Template Render & Write Helper
**Goal:** Replace verbose render-then-write pattern (15+ instances) with single helper function; ~50 lines reduced.

**Location:** Add function to `internal/common/generator.go` (created in Phase 2.3)

**Tasks:**
1. Add to `internal/common/generator.go`:
   ```go
   // RenderAndWriteTemplate renders a template and writes to file in one operation.
   func RenderAndWriteTemplate(embedFS embed.FS, outputDir, templateName, outputPath string, data any) error {
       content, err := fs.ReadFile(embedFS, filepath.Join("templates", templateName))
       if err != nil {
           return fmt.Errorf("read template %s: %w", templateName, err)
       }
       
       rendered, err := RenderTemplate(templateName, string(content), data)
       if err != nil {
           return err
       }
       
       fullPath := filepath.Join(outputDir, outputPath)
       if err := WriteFile(fullPath, rendered); err != nil {
           return err
       }
       
       return nil
   }
   ```

2. Update `internal/providers/android/generator.go`:
   - Before:
     ```go
     s, err := renderTemplate("AndroidManifest.xml.tmpl", data)
     if err != nil {
         return err
     }
     if err := writeFile(filepath.Join(outputDir, "AndroidManifest.xml"), s); err != nil {
         return err
     }
     ```
   - After:
     ```go
     if err := common.RenderAndWriteTemplate(templateFS, outputDir, "AndroidManifest.xml.tmpl", "AndroidManifest.xml", data); err != nil {
         return err
     }
     ```
   - Replace all 15+ instances in GenerateProject()

3. Update `internal/providers/springboot/generator.go` similarly

4. Update unit tests

5. Run tests: `go test ./...`

**Verification:**
```bash
go test ./...
# Should see significantly fewer lines with renderTemplate/writeFile patterns
```

**Estimated Effort:** 45 minutes
**Affected Files:** android/generator.go, springboot/generator.go, common/generator.go

---

### 3.3 Break Down Large GenerateProject Functions
**Goal:** Reduce GenerateProject from 330 lines (Android) and 230 lines (Spring Boot) to focused <100-line functions.

**Android Provider (internal/providers/android/generator.go):**

**Tasks:**
1. Refactor GenerateProject to 20-30 lines calling smaller helpers:
   ```go
   func GenerateProject(cfg providers.ProviderConfig) error {
       params, err := a.extractAndValidateConfig(cfg)
       if err != nil {
           return err
       }
       
       if err := a.generateProjectStructure(params); err != nil {
           return err
       }
       
       if err := a.generateSourceFiles(params); err != nil {
           return err
       }
       
       if err := a.generateResourceFiles(params); err != nil {
           return err
       }
       
       if err := a.setupGradleBuild(params); err != nil {
           return err
       }
       
       return a.postGenerate(params)
   }
   
   // Helper functions (each <80 lines)
   func (a *AndroidProvider) extractAndValidateConfig(cfg ProviderConfig) (*ProjectParams, error) { ... }
   func (a *AndroidProvider) generateProjectStructure(params *ProjectParams) error { ... }
   func (a *AndroidProvider) generateSourceFiles(params *ProjectParams) error { ... }
   func (a *AndroidProvider) generateResourceFiles(params *ProjectParams) error { ... }
   func (a *AndroidProvider) setupGradleBuild(params *ProjectParams) error { ... }
   func (a *AndroidProvider) postGenerate(params *ProjectParams) error { ... }
   ```

2. Move related logic into appropriate helper functions

3. Create `internal/providers/android/generator_helpers.go` for extracted helpers

4. Add unit tests for each helper function

5. Run tests: `go test ./...`

**Estimated Effort:** 1.5 hours
**Affected Files:** android/generator.go, android/generator_helpers.go (new), android/generator_test.go

---

**Spring Boot Provider (internal/providers/springboot/generator.go):**

Apply similar refactoring; structure:
```go
func GenerateProject(cfg ProviderConfig) error {
   params, err := s.extractAndValidateConfig(cfg)
   // ...
   generateBuildFiles()
   generateSourceFiles()
   generateResourceFiles()
   postGenerate()
}
```

**Tasks:**
1. Refactor similarly to Android
2. Create `internal/providers/springboot/generator_helpers.go`
3. Add unit tests
4. Run tests: `go test ./...`

**Estimated Effort:** 1 hour
**Affected Files:** springboot/generator.go, springboot/generator_helpers.go (new)

---

### 3.4 Reduce Nesting in Add Function
**Goal:** Break down Add function (280+ lines, 4-5 nesting levels) into focused step functions.

**Location:** `internal/providers/android/add.go`

**Tasks:**
1. Refactor Add to use step pattern:
   ```go
   func (a *AndroidProvider) Add(cfg ProviderConfig) error {
       input, err := a.parseAddInput(cfg)
       if err != nil {
           return err
       }
       
       steps := []func(*AddInput) error{
           a.generateResourceFile,
           a.generateLayoutIfNeeded,
           a.generateViewModelIfRequested,
           a.setupDependencyInjection,
           a.setupNavigation,
       }
       
       for _, step := range steps {
           if err := step(input); err != nil {
               return err
           }
       }
       
       return nil
   }
   ```

2. Extract each major step into its own function:
   ```go
   func (a *AndroidProvider) generateResourceFile(input *AddInput) error { ... }
   func (a *AndroidProvider) generateLayoutIfNeeded(input *AddInput) error { ... }
   func (a *AndroidProvider) generateViewModelIfRequested(input *AddInput) error { ... }
   func (a *AndroidProvider) setupDependencyInjection(input *AddInput) error { ... }
   func (a *AndroidProvider) setupNavigation(input *AddInput) error { ... }
   ```

3. Move extracted functions to `internal/providers/android/add_helpers.go`

4. Update all unit tests to test individual steps

5. Run tests: `go test ./...`

**Estimated Effort:** 1.5 hours
**Affected Files:** android/add.go, android/add_helpers.go (new), android/add_test.go (update)

---

## Phase 4: Testing & Architecture Improvements (2-3 hours)

These refactorings improve testability and architectural clarity.

### 4.1 Create Test Helpers & Builder Pattern
**Goal:** Reduce test setup boilerplate by 90%; ~40 lines → 3 lines per test.

**Files to create/update:**
- `internal/providers/android/testhelpers.go` (new)
- Update: `internal/providers/android/add_test.go`, `add_flags_test.go`, `add_di_nav_test.go`, `add_edge_cases_test.go`

**Tasks:**
1. Create `internal/providers/android/testhelpers.go`:
   ```go
   package android
   
   import (
       "path/filepath"
       "testing"
       "github.com/backsofangels/grimoire/internal/providers"
   )
   
   // AddTestBuilder provides a fluent interface for building test configs.
   type AddTestBuilder struct {
       t      *testing.T
       tmp    string
       module string
       cfg    providers.ProviderConfig
   }
   
   // NewAddTestBuilder creates a new test builder with defaults.
   func NewAddTestBuilder(t *testing.T) *AddTestBuilder {
       tmp := t.TempDir()
       return &AddTestBuilder{
           t:      t,
           tmp:    tmp,
           module: filepath.Join(tmp, "app"),
           cfg: providers.ProviderConfig{
               "Module":      filepath.Join(tmp, "app"),
               "Lang":        "kotlin",
               "Kind":        "activity",
               "PackageName": "com.example.test",
           },
       }
   }
   
   func (b *AddTestBuilder) WithKind(kind string) *AddTestBuilder {
       b.cfg["Kind"] = kind
       return b
   }
   
   func (b *AddTestBuilder) WithName(name string) *AddTestBuilder {
       b.cfg["Name"] = name
       return b
   }
   
   func (b *AddTestBuilder) WithLang(lang string) *AddTestBuilder {
       b.cfg["Lang"] = lang
       return b
   }
   
   func (b *AddTestBuilder) WithDI(di string) *AddTestBuilder {
       b.cfg["DI"] = di
       return b
   }
   
   func (b *AddTestBuilder) WithUI(ui string) *AddTestBuilder {
       b.cfg["UI"] = ui
       return b
   }
   
   func (b *AddTestBuilder) WithOverride(v bool) *AddTestBuilder {
       b.cfg["Override"] = v
       return b
   }
   
   func (b *AddTestBuilder) Build() (string, providers.ProviderConfig) {
       return b.module, b.cfg
   }
   ```

2. Update test files to use builder:
   - Before:
     ```go
     func TestAddActivityCreatesFiles(t *testing.T) {
         tmp := t.TempDir()
         module := filepath.Join(tmp, "app")
         cfg := providers.ProviderConfig{
             "Kind":        "activity",
             "Name":        "MyActivity",
             "PackageName": "com.example.test",
             "Module":      module,
             "Lang":        "kotlin",
         }
         // ... test
     }
     ```
   - After:
     ```go
     func TestAddActivityCreatesFiles(t *testing.T) {
         module, cfg := NewAddTestBuilder(t).
             WithKind("activity").
             WithName("MyActivity").
             Build()
         // ... test
     }
     ```

3. Apply to all add_test.go files: reduce ~40 → 3 lines per test setup

4. Run tests: `go test ./...`

**Estimated Effort:** 1 hour
**Affected Files:** android/testhelpers.go (new), add_test.go, add_flags_test.go, add_di_nav_test.go, add_edge_cases_test.go

---

### 4.2 Parameterize Validation Tests
**Goal:** Reduce 3 identical test functions (50 lines each) to 1 parameterized test (20 lines).

**Location:** `cmd/add_validation_test.go`

**Tasks:**
1. Create generic table-driven helper (or use subtests):
   ```go
   func TestValidators(t *testing.T) {
       tests := []struct {
           name      string
           validator func(string) error
           cases     []struct {
               input   string
               wantErr bool
           }
       }{
           {
               name:      "Language",
               validator: validator.ValidateLanguage,
               cases: []struct {
                   input   string
                   wantErr bool
               }{
                   {"", false},
                   {"kotlin", false},
                   {"java", false},
                   {"swift", true},
               },
           },
           {
               name:      "UI",
               validator: validator.ValidateUI,
               cases: []struct {
                   input   string
                   wantErr bool
               }{
                   {"", false},
                   {"xml", false},
                   {"compose", false},
                   {"none", false},
                   {"invalid", true},
               },
           },
           {
               name:      "DI",
               validator: validator.ValidateDI,
               cases: []struct {
                   input   string
                   wantErr bool
               }{
                   {"", false},
                   {"none", false},
                   {"hilt", false},
                   {"koin", false},
                   {"dagger", true},
               },
           },
       }
       
       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               for _, c := range tt.cases {
                   if err := tt.validator(c.input); (err != nil) != c.wantErr {
                       t.Errorf("validator(%q) wantErr=%v gotErr=%v", c.input, c.wantErr, err)
                   }
               }
           })
       }
   }
   ```

2. Remove old TestValidateUI, TestValidateLang, TestValidateDI functions

3. Run tests: `go test ./...`

**Verification:**
```bash
go test -v cmd/... | grep -i validation
# Should see one test with subtests
```

**Estimated Effort:** 30 minutes
**Affected Files:** cmd/add_validation_test.go

---

### 4.3 Create Base Provider Implementation
**Goal:** Reduce boilerplate in each provider; shared default implementations for Add() and DoctorChecks().

**Location:** `internal/providers/provider.go`

**Tasks:**
1. Add BaseProvider struct to `internal/providers/provider.go`:
   ```go
   // BaseProvider provides default implementations of Provider methods.
   type BaseProvider struct {
       name        string
       description string
   }
   
   func NewBaseProvider(name, description string) *BaseProvider {
       return &BaseProvider{
           name:        name,
           description: description,
       }
   }
   
   func (b *BaseProvider) Name() string {
       return b.name
   }
   
   func (b *BaseProvider) Description() string {
       return b.description
   }
   
   func (b *BaseProvider) Add(cfg ProviderConfig) error {
       return fmt.Errorf("add: not supported for %s provider", b.name)
   }
   
   func (b *BaseProvider) DoctorChecks() []Check {
       return nil
   }
   ```

2. Update providers to embed BaseProvider:
   - `internal/providers/springboot/provider.go`:
     ```go
     type SpringBootProvider struct {
         *providers.BaseProvider
     }
     
     func New() Provider {
         return &SpringBootProvider{
             BaseProvider: providers.NewBaseProvider(
                 "springboot",
                 "Spring Boot project provider (Java)",
             ),
         }
     }
     
     // Only implement: Flags, Prompt, Validate, Generate
     // Add and DoctorChecks inherited from base
     ```

3. Remove duplicate Name() and Description() implementations

4. Update Android provider similarly (but keep custom DoctorChecks)

5. Run tests: `go test ./...`

**Verification:**
```bash
go test ./...
# Verify Add() returns proper error for springboot
```

**Estimated Effort:** 45 minutes
**Affected Files:** internal/providers/provider.go, android/provider.go, springboot/provider.go

---

### 4.4 Enhanced Config Loading (internal/config/config.go)
**Goal:** Add validation and status tracking to config loading; better error diagnostics.

**Tasks:**
1. Define ConfigLoadResult in `internal/config/config.go`:
   ```go
   type ConfigLoadResult struct {
       Config    Config
       Loaded    bool      // Was config actually loaded or defaults used?
       FilePath  string
       Warnings  []string
   }
   ```

2. Add validateConfig function:
   ```go
   func validateConfig(c Config) []string {
       var warnings []string
       if c.DefaultMinSdk < 21 {
           warnings = append(warnings, "default_min_sdk < 21 (minimum supported is 21)")
       }
       if c.DefaultMinSdk > c.DefaultTargetSdk {
           warnings = append(warnings, "default_min_sdk > default_target_sdk")
       }
       if c.DefaultLang != "kotlin" && c.DefaultLang != "java" && c.DefaultLang != "" {
           warnings = append(warnings, fmt.Sprintf("invalid default_lang: %s", c.DefaultLang))
       }
       return warnings
   }
   ```

3. Refactor Load() to return ConfigLoadResult:
   ```go
   func Load() (ConfigLoadResult, error) {
       cfg := DefaultConfig()
       result := ConfigLoadResult{Config: cfg, Loaded: false}
       
       home, err := os.UserHomeDir()
       if err != nil {
           result.Warnings = append(result.Warnings, "Could not determine home directory")
           return result, nil
       }
       
       path := filepath.Join(home, ".grimoire", "config.json")
       result.FilePath = path
       
       data, err := os.ReadFile(path)
       if os.IsNotExist(err) {
           return result, nil  // Defaults used, not loaded
       }
       if err != nil {
           return result, fmt.Errorf("read config: %w", err)
       }
       
       if err := json.Unmarshal(data, &cfg); err != nil {
           return result, fmt.Errorf("parse config: %w", err)
       }
       
       result.Config = cfg
       result.Loaded = true
       result.Warnings = validateConfig(cfg)
       
       return result, nil
   }
   ```

4. Update cmd/root.go to handle warnings:
   ```go
   result, err := config.Load()
   if err != nil {
       logging.Error("Failed to load config", "error", err)
       return
   }
   for _, w := range result.Warnings {
       logging.Warn("Config warning", "message", w)
   }
   ```

5. Add unit tests for validateConfig

6. Run tests: `go test ./...`

**Estimated Effort:** 45 minutes
**Affected Files:** internal/config/config.go, cmd/root.go, config tests

---

## Testing & Verification Strategy

### Per-Phase Testing
After each phase, run:
```bash
go build ./...          # Verify compilation
go test ./...           # Run all tests
go test -race ./...     # Check for race conditions
gofmt -s ./...          # Verify formatting
go vet ./...            # Check for issues
```

### Integration Testing
After Phases 3-4, run manual smoke tests:
```bash
go build -o grimoire.exe .
./grimoire.exe new TestApp --lang kotlin
./grimoire.exe new TestApp --lang java --template compose  # Should error
./grimoire.exe add activity MyActivity
./grimoire.exe doctor
```

### Test Coverage
Maintain >80% test coverage:
```bash
go test -cover ./...
```

---

## Rollback Strategy

Each phase is independently reversible:
- **Phase 1**: Cosmetic changes; easy to revert with `git diff`
- **Phase 2**: New helper packages; revert by removing imports and reverting to old patterns
- **Phase 3**: Function refactoring; revert by removing helper functions and inlining logic
- **Phase 4**: Testing/architecture; revert by removing base providers and test builders

**Recommended:** Commit after each refactoring task; use `git revert` if needed.

---

## Success Criteria

| Metric | Target | Phase |
|--------|--------|-------|
| Duplicate code eliminated | 300+ lines | 3 |
| Boilerplate reduced | 200+ lines | 2-3 |
| Test setup reduction | 90% | 4 |
| Function avg size | <100 lines | 3 |
| Test coverage maintained | >80% | All |
| All tests passing | 100% | All |
| No breaking changes | Yes | All |

---

## Timeline & Resource Allocation

| Phase | Tasks | Est. Time | Dependencies | Can Parallelize? |
|-------|-------|-----------|--------------|-----------------|
| 1 | 4 | 1-2h | None | Partial (1.3 depends on 1.1-1.2) |
| 2 | 3 | 2-3h | Phase 1 | No (2.1 and 2.2 independent; 2.3 independent) |
| 3 | 4 | 3-4h | Phase 2 | Partial (3.1, 3.2 independent; 3.3, 3.4 independent) |
| 4 | 4 | 2-3h | Phase 3 | Partial (4.1, 4.2 independent; 4.3, 4.4 independent) |
| **Total** | **15** | **~10-14h** | Sequential | **Yes, with planning** |

**Recommended staffing:** 1 developer, 10-14 hours over 2-3 working days (or 3-4 half-days with short breaks).

---

## Post-Refactoring Tasks

1. **Update README.md** with new architecture (if applicable)
2. **Update AGENTS.md** with completed refactorings
3. **Update EXAMPLES.md** if CLI patterns changed
4. **Create CHANGELOG entry** documenting improvements
5. **Review test coverage** and add new tests as needed
6. **Performance benchmarking** (optional) to ensure no regressions
7. **Documentation review** of new helper packages

---

## Related Documentation

- See [AGENTS.md](AGENTS.md) for architecture rules and patterns
- See refactoring analysis for detailed problem descriptions and code examples
- See [README.md](README.md) for current feature set and usage

