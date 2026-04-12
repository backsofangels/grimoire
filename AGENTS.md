# AGENTS.md — Grimoire 🔮

This document provides instructions and context for AI coding agents working on the Grimoire codebase.
Read this file entirely before making any changes.

---

## Project Summary

**Grimoire** is a CLI scaffolding tool written in **Go**, designed to initialize development projects from the terminal without requiring an IDE. The primary (and currently only) provider is **Android**, targeting **VSCode** as the development environment.

The tool is inspired by `flutter create` and `ng new`: it supports both an interactive wizard mode and a fully non-interactive flag-based mode.

Refer to `README.md` for full feature documentation, command reference, and project structure.

---

## Agent Progress

- **Completed:** Repo bootstrap; Provider interface and registry; Android provider skeleton; Validator module and tests; Templates embedding; Generator implementation; `grimoire new` command (interactive TUI); Generator unit tests; Centralized logging wrapper (`internal/logging`) with branded helpers (`Success`, `Step`); Translation of user-facing strings to English; `grimoire add` command and provider-level `Add` API (activity|fragment|viewmodel) with DI (Hilt/Koin) and navigation wiring; Jetpack Compose template and Compose-ready Gradle support; Many unit tests covering add/generator flows (Kotlin & Java) — tests passing locally; Removed embedded Gradle wrapper assets; Runtime wrapper strategy implemented (uses system `gradle` and `--no-wrapper` option); Cleaned temporary `gradle-tmp` and updated `.gitignore`; README updates documenting `add` and Compose templates.
- **In-Progress:** Java availability and network checks in the generator; per-OS Gradle install instructions; documentation polish and manual QA / smoke builds across platforms.
- **Pending:** Release configuration and CI pipeline; Scoop distribution and GoReleaser finalization; cross-platform verification and packaging.

**Smoke Test Fixes Applied:**

- Added Kotlin Gradle plugin classpath (`org.jetbrains.kotlin:kotlin-gradle-plugin:1.9.23`) to top-level `build.gradle` template.
- Added `org.jetbrains.kotlin:kotlin-stdlib:1.9.23` to app dependencies for Kotlin projects.
- Added `namespace` to the module `build.gradle` template and removed `package` from generated `AndroidManifest.xml` (AGP expects `namespace` in build files).
- Added Android XML namespace (`xmlns:android`) and `android:exported="true"` to the activity manifest entry (required for Android 12+).
- Ensured `compileOptions` and `kotlinOptions` align (Java/Kotlin target set to 1.8) to avoid JVM target mismatches on machines running newer JDKs.

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
│   ├── providers/
│   │   ├── provider.go           # Provider interface + ProviderConfig + Check types
│   │   ├── registry.go           # Global provider registry (register all providers here)
│   │   └── android/              # Android provider (reference implementation)
│   │       ├── provider.go       # Implements Provider interface
│   │       ├── generator.go      # File generation logic
│   │       ├── doctor.go         # Environment checks (JDK, ANDROID_HOME, etc.)
│   │       ├── prompts.go        # Interactive wizard using charmbracelet/huh
│   │       └── templates/        # Embedded Go templates (.tmpl)
│   │           ├── AndroidManifest.xml.tmpl
│   │           ├── MainActivity.kt.tmpl
│   │           ├── MainActivity.java.tmpl
│   │           ├── build_gradle_app.tmpl
│   │           ├── build_gradle_root.tmpl
│   │           ├── settings_gradle.tmpl
│   │           ├── gradle_properties.tmpl
│   │           ├── activity_main_xml.tmpl
│   │           ├── strings_xml.tmpl
│   │           ├── themes_xml.tmpl
│   │           ├── gitignore.tmpl
│   │           ├── vscode_settings.tmpl
│   │           └── vscode_extensions.tmpl
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
    - `--ui` / `--no-ui`: allowed values `xml | compose | none`. `--no-ui` is equivalent to `--ui none`.
    - `--lang`: allowed values `kotlin | java`.
    - `--di`: allowed values `none | hilt | koin`.
- **Interactive TUI:** prefill values from flags/config, but still enforce the same validations. If a flag is invalid, the CLI should display a user-friendly error instead of running generation.
- **Implementation note:** Add small helper functions in `cmd/` (for example `validateUI`, `validateLang`, `validateDI`) and call them from both the non-interactive subcommands and the interactive `runAddInteractive` flow.

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

- [ ] `Provider` interface in `internal/providers/provider.go`
- [ ] Provider registry in `internal/providers/registry.go`
- [ ] Android provider skeleton in `internal/providers/android/provider.go`
- [ ] `grimoire new` command (non-interactive, all flags explicit)
- [ ] Templates: `empty` and `basic` (Kotlin only)
- [ ] `.vscode/settings.json` and `.vscode/extensions.json` generation
- [ ] Git init via `os/exec`
- [ ] Gradle wrapper files included as embedded assets
- [ ] `validator.go`: package name and app name validation

### v0.2.0 — Environment checks

- [ ] `grimoire doctor` command
- [ ] Android doctor checks: JDK, JAVA_HOME, ANDROID_HOME, Build-Tools, Platform
- [ ] `--fix` flag: auto-set JAVA_HOME on Windows via `setx`

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

- [ ] `grimoire config set/get/list/reset`
- [ ] `~/.grimoire/config.json` read/write via `internal/config`

### v1.0.0 — Stable + Scoop distribution

- [x] Java language support in Android provider (`--lang java`)
- [ ] GoReleaser config finalized (Windows amd64 + arm64)
- [ ] GitHub Actions release workflow
- [ ] Scoop manifest in separate `scoop-bucket` repository
- [ ] Full test coverage for generator and validator

### v1.x — Future providers

- [ ] Spring Boot provider (`--provider springboot`)
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
