# goimportmaps

> Visualize and validate package dependencies in your Go project.
>

![Screenshot](./assets/html_report.png)

## Overview

`goimportmaps` is a CLI tool that analyzes your Go project's internal package imports, visualizes them as a dependency
graph, and detects architectural violations like forbidden dependencies or unexpected imports.

Ideal for large-scale or monorepo Go applications, this tool helps ensure architectural integrity by preventing
undesired package-level coupling.

## Features

- 📊 Visualize internal package dependencies (Mermaid, Graphviz, HTML)
- 🚨 Detect forbidden imports based on custom rules (`forbidden` mode)
- 🛡 Enforce allowed imports strictly (`allowed` mode, whitelist style)
- 📈 **Coupling metrics analysis** (inspired by NDepend)
  - **Afferent Coupling (Ca)**: Number of packages depending on this package
  - **Efferent Coupling (Ce)**: Number of packages this package depends on
  - **Instability (I)**: Ce / (Ca + Ce) ratio (0=stable, 1=unstable)
- ✅ Output violations with actionable messages
- 🔍 Highlight architectural drift in pull requests
- 🧠 Perfect for layered, hexagonal, or clean architecture

## Installation

```bash
# Install goimportmaps into your project
go get -tool github.com/mickamy/goimportmaps/cmd/goimportmaps@latest

# or install it globally
go install github.com/mickamy/goimportmaps/cmd/goimportmaps@latest
```

## Usage

```bash
goimportmaps ./...
```

You can also scope the analysis to specific subdirectories:

```bash
goimportmaps ./internal/...
```

### Options

| Option      | Description                                             |
|-------------|---------------------------------------------------------|
| `--format`  | Output format: `text`, `mermaid`, `html`, or `graphviz` |
| `--mode`    | Validation mode: `forbidden` (default) or `allowed`     |
| `--metrics` | Show coupling metrics (overrides config setting)       |

## Example

Given the following structure:

```
main
├── handler
│   └── user_handler.go (imports service)
├── service
│   └── user_service.go (imports infra)
└── infra
    └── db.go
```

### Forbidden Mode (default)

If `handler` imports `infra` directly, the tool will detect:

```bash
🚨 1 violation(s) found

🚨 Violation: github.com/your/project/internal/handler imports github.com/your/project/internal/infra
```

### Allowed Mode (strict whitelist)

You can enforce exact allowed imports:

```yaml
allowed:
  - source: github.com/your/project/internal/handler
    imports:
      - github.com/your/project/internal/service
    stdlib: true # allow importing standard library packages as well (default: true)
```

If `handler` tries to import anything other than `service` or stdlib (e.g., `infra`), it will be flagged.

---

### Mermaid Output Example

```markdown
graph TD
  main --> handler
  handler --> service
  service --> infra
  handler --> infra %% ❌
```

---

## Coupling Metrics

Display package coupling metrics to identify architectural issues:

```bash
# Show metrics in CLI output
goimportmaps ./... --metrics

# Generate HTML report with metrics dashboard
goimportmaps ./... --format=html --metrics > metrics-report.html
```

### Metrics Output Example

```
📊 Coupling Metrics:

Package                           Ca   Ce      I Status
internal/cli                       0   11   1.00 🚨
internal/prints                    2   11   0.85 🚨  
internal/config                    4    6   0.60 ✅

🚨 Coupling Violations:
- internal/cli: High efferent coupling (11 > 10), High instability (1.00 > 0.80)
- internal/prints: High efferent coupling (11 > 10), High instability (0.85 > 0.80)
```

### Understanding Metrics

- **Ca (Afferent Coupling)**: Higher values indicate packages that are heavily depended upon (stable)
- **Ce (Efferent Coupling)**: Higher values indicate packages with many dependencies (unstable)
- **I (Instability)**: 
  - **0.0** = Completely stable (only used by others)
  - **1.0** = Completely unstable (only uses others)
  - **0.5** = Balanced coupling

---

## Configuration

`.goimportmaps.yaml`

### Forbidden Mode Example

```yaml
forbidden:
  - source: github.com/your/project/internal/handler
    imports:
      - github.com/your/project/internal/infra
  - source: github.com/your/project/internal/app
    imports:
      - github.com/your/project/internal/db
```

### Allowed Mode Example

```yaml
allowed:
  - source: github.com/your/project/internal/handler
    imports:
      - github.com/your/project/internal/service
  - source: github.com/your/project/internal/service
    imports:
      - github.com/your/project/internal/repository
    stdlib: false
```

### Metrics Configuration Example

```yaml
metrics:
  enabled: true
  coupling:
    max_efferent: 10      # Maximum efferent coupling (Ce)
    max_afferent: 15      # Maximum afferent coupling (Ca)  
    max_instability: 0.8  # Maximum instability (I)
    warn_efferent: 7      # Warning threshold for Ce
    warn_afferent: 10     # Warning threshold for Ca
    warn_instability: 0.6 # Warning threshold for I
```

## HTML Output

Use `--format=html` to generate a standalone static report:

```bash
# Basic HTML report
goimportmaps ./... --format=html > report.html

# HTML report with metrics dashboard
goimportmaps ./... --format=html --metrics > metrics-report.html
```

The metrics-enabled HTML report includes:
- Interactive dependency graph with Mermaid.js
- Coupling metrics table with color-coded violations
- Visual instability bars and status indicators
- Summary cards showing violation counts

Reports can be viewed in your browser or uploaded as CI artifacts.

## License

[MIT](./LICENSE)
