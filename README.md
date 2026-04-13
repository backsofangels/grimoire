# 🔮 Grimoire

> Scaffold development projects from your terminal — no IDE required.

Grimoire is a lightweight CLI tool that generates production-ready projects in seconds.
Inspired by `flutter create` and `ng new`, it supports both a fully flag-based non-interactive
mode and an interactive TUI wizard with consistent, branded styling.

**Supported targets:** Android, Spring Boot (Gradle/Maven)

```bash
grimoire new MyApp --package com.example.myapp --lang kotlin --provider android
```

---

## Features

- **Multi-provider architecture** — pluggable scaffolding targets (`android`, `springboot`)
- **Android:** Kotlin & Java, full template support (basic, compose, empty)
- **Spring Boot:** Java 17+, Gradle or Maven, multiple templates
- **Gradle wrapper** — reproducible builds included by default; optional system gradle fallback
- **Beautiful TUI** — interactive wizard with consistent purple/violet branding
- **Config persistence** — saved defaults in `~/.grimoire/config.json`
- **Opinionated defaults** — AndroidManifest, VSCode config, `.gitignore`, themes, strings
- **Validation built-in** — enforces valid package names and app names before writing files
- **Environment checks** — `grimoire doctor` detects missing tools
- **CI-friendly** — fully scriptable, non-interactive mode for automation

---

## Prerequisites

| Tool | Version | Notes |
| --- | --- | --- |
| Go | 1.22+ | To build from source |
| Java JDK | 11+ (Android), 17+ (Spring Boot) | `java` must be on PATH |
| Gradle CLI | any | Optional; used to generate wrapper if present |
| Android SDK | — | Required only to run builds locally |
| Git | — | Optional; used for `git init` during creation |

---

## Installation

### Build from source

```bash
git clone https://github.com/backsofangels/grimoire
cd grimoire
go build -o grimoire.exe .
```

### Windows (Scoop)

Scoop distribution available in v1.0.0 release.

---

## Usage

### Interactive Wizard (Recommended)

```bash
grimoire new
```

Step through a beautiful colored TUI to configure your project:

1. App name (e.g., MyApp)
2. Package name (e.g., com.example.myapp)
3. Output directory
4. Language (Kotlin/Java)
5. Minimum SDK (Android)
6. Template (basic/compose/empty)
7. **Gradle wrapper** (use gradlew or system gradle)
8. Initialize git?
9. Generate VSCode config?
10. Confirm and create

### Non-interactive Mode (Scripts/CI)

```bash
# Create Kotlin Android project
grimoire new MyApp \
  --package com.example.myapp \
  --lang kotlin \
  --template basic \
  --provider android

# Create Java Spring Boot project with Maven
grimoire new MyJavaApp \
  --package com.example.javaapp \
  --provider springboot \
  --build-system maven

# Skip Gradle wrapper
grimoire new MyApp --package com.example.myapp --wrapper=false
```

### Build & Run

```bash
cd MyApp
./gradlew assembleDebug  # If wrapper enabled
# or
gradle assembleDebug     # If using system gradle
```

### Add Components to Existing Project

```bash
# Add Activity with Hilt DI and ViewModel
grimoire add activity --name MyScreen --di hilt --viewmodel --nav

# Add Fragment with Koin
grimoire add fragment --name MyFragment --di koin --viewmodel --nav

# Add only a ViewModel
grimoire add viewmodel --name MyViewModel
```

Supports:

- `--di none | hilt | koin` — dependency injection wiring
- `--viewmodel` — generate accompanying ViewModel
- `--nav` — add navigation graph entry
- `--override` — replace existing files

### Environment Checks

```bash
grimoire doctor --fix
```

Validates JDK, Gradle, Android SDK, and can auto-configure JAVA_HOME on Windows.

---

## Templates

### Android

| Name | Description |
| --- | --- |
| `basic` | Activity + layout XML + ViewModel, standard project structure |
| `empty` | Bare MainActivity, no layout directory |
| `compose` | Jetpack Compose UI (Kotlin only), Compose-ready Gradle config |

### Spring Boot

| Name | Description |
| --- | --- |
| `springboot` | Spring Boot application with starter dependencies |
| `plain` | Plain Java application |

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
- [x] Scoop distribution (v1.0.0)

---

## Contributing

1. Fork the repo and create a feature branch
2. Follow Go 1.22+ conventions — run `gofmt -w .` before committing
3. Add tests for new functionality and verify with `go test ./...`
4. Open a pull request with a clear description

---

## License

Apache 2.0 — see [`LICENSE`](./LICENSE)
