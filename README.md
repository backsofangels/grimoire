# 🔮 Grimoire

> Scaffold Android projects from your terminal — no IDE required.

Grimoire is a lightweight CLI tool that generates production-ready Android projects in seconds.
Inspired by `flutter create` and `ng new`, it supports both a fully flag-based non-interactive
mode and an interactive TUI wizard.

```bash
grimoire new MyApp --package com.example.myapp --lang kotlin
```

---

## Features

- **Provider-based architecture** — pluggable scaffolding targets (currently `android`, `springboot`)
- **Kotlin & Java** — full template support for both languages
- **Gradle-ready** — generates `build.gradle`, `settings.gradle`, `gradle.properties` and a working wrapper
- **Opinionated defaults** — `AndroidManifest.xml`, `.vscode/` config, `.gitignore`, themes, strings, layouts
- **Validation built-in** — enforces valid package names and app names before writing a single file
- **CI-friendly** — fully scriptable, no prompts unless you want them

---

## Prerequisites

| Tool | Version | Notes |
| --- | --- | --- |
| Go | 1.22+ | To build from source |
| Java JDK | 11+ | `java` must be on PATH |
| Gradle CLI | any | Used to generate the wrapper; install hints shown if missing |
| Android SDK | — | Required only to run `assembleDebug` locally |
| Git | — | Recommended; used for `git init` during project creation |

---

## Installation

### Build from source

```bash
git clone https://github.com/backsofangels/grimoire
cd grimoire
go build -o grimoire.exe .
```

### Windows (Scoop) — coming in v1.0.0

```powershell
scoop bucket add grimoire https://github.com/backsofangels/scoop-bucket
scoop install grimoire
```

---

## Usage

### Create a Kotlin project

```bash
grimoire new MyApp --package com.example.myapp --lang kotlin --template basic
```

### Create a Java project

```bash
grimoire new MyJavaApp --package com.example.javaapp --lang java --template basic
```

### Interactive wizard (no flags needed)

```bash
grimoire new
```

The wizard guides you step by step through app name, package, language, SDK version, template, and extras.

### Skip Gradle wrapper generation

```bash
grimoire new MyApp --package com.example.myapp --wrapper=false
```

### Build the generated project

```bash
cd MyApp
./gradlew assembleDebug
```

> Requires the Android SDK and matching build tools configured locally.

### Add subcommands

Use `grimoire add` to generate UI components and their supporting files inside an existing project.

```bash
# Add an Activity with Hilt, ViewModel, and nav entry
grimoire add activity --name MyActivity --di hilt --viewmodel --nav

# Add a Fragment with Koin and nav entry
grimoire add fragment --name MyFragment --di koin --viewmodel --nav

# Create only a ViewModel
grimoire add viewmodel --name MyViewModel --di hilt
```

- `--name`: explicit class/resource name (defaults derived from the command)
- `--di`: dependency injection wiring; one of `none`, `hilt`, or `koin` (default: `none`)
- `--viewmodel`: also generate a ViewModel and wire it to the UI component
- `--nav`: add an entry to `app/src/main/res/navigation/nav_graph.xml`
- `--override`: overwrite existing files when present

- Note: omitting `--di` (or using `--di none`) generates the component without any DI wiring.

---

## Templates

| Name | Description |
| --- | --- |
| `basic` | Activity + layout XML + ViewModel, full project structure |
| `empty` | Bare `MainActivity`, no layout directory |
| `compose` | Jetpack Compose starter template (Kotlin only) with Compose-ready Gradle config |

---

## Project Structure

Generated projects follow standard Android conventions and open immediately in VSCode or Android Studio:

```text
MyApp/
├── app/
│ ├── src/main/
│ │ ├── java/com/example/myapp/
│ │ │ └── MainActivity.kt
│ │ ├── res/
│ │ │ ├── layout/activity_main.xml
│ │ │ ├── values/strings.xml
│ │ │ └── values/themes.xml
│ │ └── AndroidManifest.xml
│ └── build.gradle
├── build.gradle
├── settings.gradle
├── gradle.properties
├── .vscode/
│ ├── settings.json
│ └── extensions.json
└── .gitignore
```

---

## Development

### Testing

- Run unit and fast tests:

```bash
go test ./...
```

- Integration tests (resource-heavy):
  - Integration tests are under the `cmd/` package and exercise the interactive TUI and end-to-end generator flow.
  - These tests are skipped automatically in CI: the test code checks the `CI` environment variable and calls `t.Skip()` when it is set. Most CI providers (GitHub Actions, GitLab CI, etc.) set `CI=true` by default, so heavy tests won't run in CI pipelines.
  - To run integration tests locally (they require Java and Gradle for smoke builds):

```bash
# Run a specific integration test by name
go test ./cmd -run TestNewIntegrationCLI -v
go test ./cmd -run TestCtrlCCancelFlow -v

# Or run all tests in the cmd package (including integration tests)
go test ./cmd -v
```

- Notes:
  - Ensure `java` (JDK) and `gradle` are on your `PATH` before running integration tests that build generated projects.
  - Integration tests create temporary directories and clean up after themselves; they're intended for local verification only.

Format code:

```bash
gofmt -w .
```

Templates are in `internal/providers/android/templates/` and use Go's `text/template`.
The provider interface lives in `internal/providers/provider.go` — implement it to add new targets.

---

## Roadmap

- [x] Android provider — Kotlin & Java
- [x] Gradle wrapper generation
- [x] Smoke-tested end-to-end (`assembleDebug`)
- [x] `grimoire doctor` — environment preflight checks (JDK, Gradle, Android SDK)
- [x] `grimoire new` interactive wizard (TUI)
- [x] `grimoire add activity|fragment|viewmodel`
- [x] Jetpack Compose template
- [x] Spring Boot provider
- [ ] Scoop distribution (v1.0.0)

---

## Contributing

1. Fork the repo and create a feature branch
2. Follow Go 1.22+ conventions — run `gofmt -w .` before committing
3. Add tests for new functionality and verify with `go test ./...`
4. Open a pull request with a clear description

---

## License

Apache 2.0 — see [`LICENSE`](./LICENSE)
