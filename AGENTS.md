# AGENTS.md — Grimoire 🔮

This document provides instructions and context for AI coding agents working on the Grimoire codebase.
Read this file entirely before making any changes.

---

## Project Summary

**Grimoire** is a CLI scaffolding tool written in **Go**, designed to initialize development projects from the terminal without requiring an IDE. It supports multiple providers (**Android**, **Spring Boot**), targeting **VSCode** as the primary development environment.

The tool is inspired by `flutter create` and `ng new`: it supports both an interactive TUI wizard mode and a fully non-interactive flag-based mode.

Refer to `README.md` for full feature documentation, command reference, and project structure.

---

## Agent Progress

- **Completed:** Repo bootstrap; Provider interface and registry; Android provider skeleton; Validator module and tests; Templates embedding; Generator implementation; `grimoire new` command (interactive TUI); Generator unit tests; Centralized logging wrapper (`internal/logging`) with branded helpers (`Success`, `Step`); Translation of user-facing strings to English; `grimoire add` command and provider-level `Add` API (activity|fragment|viewmodel) with DI (Hilt/Koin) and navigation wiring; Jetpack Compose template and Compose-ready Gradle support; Many unit tests covering add/generator flows (Kotlin & Java) — tests passing locally; Removed embedded Gradle wrapper assets; Runtime wrapper strategy implemented; Cleaned temporary `gradle-tmp` and updated `.gitignore`; README updates; **Phase 3 CLI flag rationalization:** replaced `--no-wrapper` with `--wrapper`; Spring Boot optional `--group` and `--artifact`; `--module` default set to "app"; `--ui` restricted to xml|compose; Java + Compose incompatibility validation; **Phase 4 Code Quality & Refactoring (Completed):** Removed unused imports and formatting cleanup; extracted string and numeric constants; created normalization utilities (`internal/common`); consolidated validation functions and created config extraction helpers; created template render helper and broke down GenerateProject functions; reduced nesting in Add function; created test helpers with builder pattern; parameterized validation tests; implemented base provider functionality; **Enhanced config system** (`internal/config/config.go`) with persistence to `~/.grimoire/config.json`, validation logic, merge semantics for CLI flags/wizard/defaults, full unit test coverage; **TUI Styling Consistency:** Created shared `internal/tui/theme.go` package with centralized color palette and formatting helpers (`RenderGroupTitle`, `PrintHeader`); refactored Android & Spring Boot prompts to eliminate ~120 lines of duplicated styling code; **Gradle wrapper TUI selection:** Step 7 wrapper choice in Android wizard with sensible defaults; **Version tag stability:** Verified AGP 8.4.0, Kotlin 1.9.22, Compose Compiler 1.5.10, Spring Boot 3.2.0 compatibility; **README.md alignment:** Updated to multi-provider documentation, added Templates table, corrected prerequisites by provider, documented wrapper selection, added Configuration/Versions/Architecture sections; all phases verified with passing test suite.
- **In-Progress:** Acceptance testing with built binary across platforms; cross-platform verification and smoke builds.
- **Pending:** `grimoire config` command (set/get/list/reset); Release configuration and CI pipeline; Scoop distribution and GoReleaser finalization.

**Smoke Test Fixes Applied:**

- Upgraded Android Gradle Plugin to `8.4.0` (from 8.3.0) for proper compileSdk 35 support without warnings.
- Updated Kotlin Gradle plugin to `1.9.22` (downgraded from 1.9.23 for Compose Compiler compatibility).
- Added `kotlin-stdlib:1.9.22` to app dependencies for Kotlin projects (compatible with Compose Compiler 1.5.10).
- Set Compose Compiler to `1.5.10` with explicit version pinning in composeOptions for Compose templates.
- Added `namespace` to the module `build.gradle` template and removed `package` from generated `AndroidManifest.xml` (AGP expects `namespace` in build files).
- Added Android XML namespace (`xmlns:android`) and `android:exported="true"` to the activity manifest entry (required for Android 12+).
- Ensured `compileOptions` and `kotlinOptions` align (Java/Kotlin target set to 1.8) to avoid JVM target mismatches.
- Added `android.suppressUnsupportedCompileSdk=<targetSdk>` to `gradle.properties` to suppress AGP warnings.
- Fixed Spring Boot **Maven** template to use Java 17 (upgraded from 11) matching Spring Boot 3.2.0 requirement.
- Added Spring Boot parent declaration to pom.xml for proper dependency management.

## Repository Layout

```text
grimoire/
├── cmd/                          # Cobra CLI commands — one file per command
│   ├── root.go                   # Root command, persistent flags, version
│   ├── new.go                    # grimoire new
│   ├── add.go                    # grimoire add
│   ├── doctor.go                 # grimoire doctor
│   └── config.go                 # grimoire config
├── internal/
│   ├── logging/                # Centralized logging wrapper (internal/logging/logging.go)
│   ├── tui/                    # Shared TUI theming (internal/tui/theme.go)
│   ├── providers/
│   │   ├── provider.go           # Provider interface + ProviderConfig + Check types
│   │   ├── registry.go           # Global provider registry (register all providers here)
│   │   ├── android/              # Android provider (reference implementation)
│   │   │   ├── provider.go       # Implements Provider interface
│   │   │   ├── generator.go      # File generation logic
│   │   │   ├── add.go            # Add activity/fragment/viewmodel logic
│   │   │   ├── doctor.go         # Environment checks (JDK, ANDROID_HOME, etc.)
│   │   │   ├── prompts.go        # Interactive wizard using charmbracelet/huh
│   │   │   └── templates/        # Embedded Go templates (.tmpl)
│   │   │       ├── AndroidManifest.xml.tmpl
│   │   │       ├── MainActivity.kt.tmpl
│   │   │       ├── MainActivity.java.tmpl
│   │   │       ├── activity.kt.tmpl
│   │   │       ├── activity.java.tmpl
│   │   │       ├── fragment.kt.tmpl
│   │   │       ├── fragment.java.tmpl
│   │   │       ├── viewmodel.kt.tmpl
│   │   │       ├── viewmodel.java.tmpl
│   │   │       ├── build_gradle_app.tmpl
│   │   │       ├── build_gradle_root.tmpl
│   │   │       ├── settings_gradle.tmpl
│   │   │       ├── gradle_properties.tmpl
│   │   │       ├── gitignore.tmpl
│   │   │       ├── vscode_settings.json.tmpl
│   │   │       └── vscode_extensions.json.tmpl
│   │   └── springboot/           # Spring Boot provider
│   │       ├── provider.go       # Implements Provider interface
│   │       ├── generator.go      # File generation logic
│   │       ├── prompts.go        # Interactive wizard
│   │       └── templates/        # Embedded Go templates (.tmpl)
│   │           ├── pom.xml.tmpl
│   │           ├── build_gradle.tmpl
│   │           ├── application.properties.tmpl
│   │           ├── gitignore.tmpl
│   │           ├── gradle_properties.tmpl
│   │           ├── settings_gradle.tmpl
│   │           ├── application_springboot.java.tmpl
│   │           └── application_plain.java.tmpl
│   ├── validator/
│   │   └── validator.go          # Shared validation: package names, app names
│   └── config/
│       └── config.go             # Read/write ~/.grimoire/config.json
├── main.go                       # Entry point — calls cmd.Execute()
├── go.mod
├── go.sum
├── .goreleaser.yaml              # GoReleaser config → Scoop distribution
├── .github/
│   └── workflows/
│       └── release.yml           # GitHub Actions — triggers GoReleaser on tag push
├── GRIMOIRE.md                   # Full project specification (source of truth)
└── AGENTS.md                     # This file
```

---

## Architecture Rules

### 1. Provider interface is the core abstraction

Every scaffolding target (Android, Spring Boot, etc.) must implement the `Provider` interface defined in `internal/providers/provider.go`. **Never add framework-specific logic to the `cmd/` layer.**

```go
type Provider interface {
    Name() string
    Description() string
    Flags() []ProviderFlag
    Prompt() (ProviderConfig, error)
    Validate(cfg ProviderConfig) error
    Generate(cfg ProviderConfig) error
    DoctorChecks() []Check
}
```

The `cmd/` layer only:

1. Resolves the active provider from `--provider` flag or `config`
2. Merges CLI flags into a `ProviderConfig` map
3. Calls the appropriate provider method

### 2. Templates use Go's `text/template`, embedded via `//go:embed`

All template files live in `internal/providers/<name>/templates/` and are embedded at compile time:

```go
//go:embed templates/*
var templateFS embed.FS
```

Use `{{.AppName}}`, `{{.PackageName}}`, `{{.MinSdk}}` etc. as template variables.
Never use string concatenation to build file contents.

### 3. Config is always read before flag parsing

The global config in `~/.grimoire/config.json` provides default values.
CLI flags override config values. Interactive wizard pre-fills from config.
Priority order: **CLI flags > wizard input > config defaults > hardcoded defaults**.

### 4. `grimoire doctor --fix` must be non-destructive

Auto-fix actions are limited to:

- Setting environment variables via `setx` (Windows) or shell profile append (Unix)
- Never modifying SDK files, PATH entries not related to Grimoire, or system settings

### 5. No third-party template engines

Use only Go's stdlib `text/template`. Do not add Handlebars, Mustache, or similar dependencies.

---

## Dependencies

| Package | Version | Purpose |
| --- | --- | --- |
| `github.com/spf13/cobra` | latest | CLI command/flag framework |
| `github.com/charmbracelet/huh` | latest | Interactive wizard UI |
| `github.com/charmbracelet/lipgloss` | latest | Terminal output styling |
| `github.com/charmbracelet/log` | latest | Styled structured logging |

Do not add dependencies without a strong justification. Prefer stdlib where possible.

---

## Coding Conventions

- **Go version**: 1.22+
- **Error handling**: always wrap errors with `fmt.Errorf("context: %w", err)`, never silently discard
- **Logging**: use `internal/logging` wrapper for all user-facing output (helpers: `Init`, `Info`, `Warn`, `Error`, `Success`, `Step`); avoid `fmt.Println`. The wrapper uses `charmbracelet/log` under the hood.
- **File paths**: always use `filepath.Join()`, never hardcode `/` or `\` separators
- **OS detection**: use `runtime.GOOS` for platform-specific logic (e.g. `setx` on Windows)
- **Tests**: every `generator.go` and `validator.go` must have a corresponding `_test.go`
- **Exported symbols**: only export what is used outside the package

### CLI Input Validation

- **Validate at the `cmd/` layer:** All user-provided CLI flag values must be validated in the `cmd/` layer before constructing a `ProviderConfig` and invoking provider methods. This prevents invalid inputs from reaching generation logic and provides fast, clear feedback to users.
- **Common validations:**
    - `--ui` / `--no-ui`: allowed values `xml | compose`. `--no-ui` disables UI generation (sets UI to "none" internally).
    - `--lang`: allowed values `kotlin | java`.
    - `--di`: allowed values `none | hilt | koin`.
    - **Android lang/template validation:** Java language is incompatible with Compose template; validation rejects `--lang java --template compose` with error message.
- **Interactive TUI:** prefill values from flags/config, but still enforce the same validations. If a flag is invalid, the CLI should display a user-friendly error instead of running generation.
- **Implementation note:** Add small helper functions in `cmd/` (for example `validateUI`, `validateLang`, `validateDI`) and call them from both the non-interactive subcommands and the interactive `runAddInteractive` flow. Conditional validation: only call `validateUI()` when UI value is not "none" (internal state set by `--no-ui`).

### Output formatting conventions

Use consistent prefixes for terminal output:

```text
✓  Success / created file
✗  Error / missing requirement
!  Warning / suggestion
→  Info / path / value
```

Colors via `lipgloss`:

- Green → success
- Red → error
- Yellow → warning
- Cyan → paths, values, highlights

- Use the `internal/logging` helpers `Success` and `Step` for multi-step and branded messages.

---

## Implementation Order (Roadmap)

Implement features in this exact order. Do not skip ahead.

### v0.1.0 — Core scaffolding

- [x] `Provider` interface in `internal/providers/provider.go`
- [x] Provider registry in `internal/providers/registry.go`
- [x] Android provider skeleton in `internal/providers/android/provider.go`
- [x] `grimoire new` command (non-interactive, all flags explicit)
- [x] Templates: `empty` and `basic` (Kotlin only)
- [x] `.vscode/settings.json` and `.vscode/extensions.json` generation
- [x] Git init via `os/exec`
- [x] Gradle wrapper files included as embedded assets
- [x] `validator.go`: package name and app name validation

### v0.2.0 — Environment checks

- [x] `grimoire doctor` command
- [x] Android doctor checks: JDK, JAVA_HOME, ANDROID_HOME, Build-Tools, Platform
- [x] `--fix` flag: auto-set JAVA_HOME on Windows via `setx`

### v0.3.0 — Interactive wizard

- [x] `grimoire new` wizard mode using `charmbracelet/huh`
- [x] Pre-fill wizard from `~/.grimoire/config.json`

### v0.4.0 — Add command

- [x] `grimoire add activity`
- [x] `grimoire add fragment`
- [x] `grimoire add viewmodel`
- [ ] `grimoire add module`

### v0.5.0 — Compose template

- [x] `--template compose` for Android provider
- [x] Compose BOM in `build.gradle`
- [x] Starter `Greeting` composable template

### v0.6.0 — Config command

- [x] `~/.grimoire/config.json` read/write persistence via `internal/config` (validation & merge logic)
- [ ] `grimoire config set/get/list/reset` command

### v1.0.0 — Stable + Scoop distribution

- [x] Java language support in Android provider (`--lang java`)
- [ ] GoReleaser config finalized (Windows amd64 + arm64)
- [ ] GitHub Actions release workflow
- [ ] Scoop manifest in separate `scoop-bucket` repository
- [x] Full test coverage for generator and validator (unit tests comprehensive; integration tests for new/add)

### v1.x — Future providers

- [x] Spring Boot provider (`--provider springboot`) — basic template implementation
- [ ] Ktor provider (`--provider ktor`)

---

## Adding a New Provider (instructions for agent)

To add a provider (e.g. `springboot`):

1. Create `internal/providers/springboot/` directory
2. Implement all methods of the `Provider` interface in `provider.go`
3. Add provider-specific flags in `Flags()` — these will be automatically wired to `grimoire new` by the registry
4. Add templates under `internal/providers/springboot/templates/`
5. Register the provider in `internal/providers/registry.go`:

```go
func init() {
    Register(android.New())
    Register(springboot.New()) // add here
}
```

6. Add doctor checks in `DoctorChecks()` (e.g. check for Maven/Gradle, JAVA_HOME)
7. Update the provider table in `README.md`

No changes to `cmd/` are needed.

---

## What NOT to do

- Do not add Android Studio project files (`.idea/`, `*.iml`)
- Do not generate `local.properties` — this is machine-specific and gitignored
- Do not hardcode SDK versions as constants in `cmd/` — they belong in each provider
- Do not use `os.Exit()` outside of `main.go` — return errors and let Cobra handle exit codes
- Do not add Windows-only code outside of `runtime.GOOS == "windows"` guards
- Do not modify `GRIMOIRE.md` or `AGENTS.md` during implementation unless explicitly instructed

---

## Key Files Reference

| File | Role |
| --- | --- |
| `AGENTS.md` | This file — instructions for coding agents |
| `internal/providers/provider.go` | Provider interface — do not change without updating all implementations |
| `internal/providers/registry.go` | Add new providers here only |
| `internal/providers/android/` | Reference implementation — follow this pattern for new providers |
| `internal/logging/logging.go` | Centralized logging wrapper used by commands and providers |
