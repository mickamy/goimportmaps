# goimportmaps

> Visualize and validate package dependencies in your Go project.

## Overview

`goimportmaps` is a CLI tool that analyzes your Go project's internal package imports, visualizes them as a dependency
graph, and detects architectural violations like forbidden dependencies (e.g., `handler -> infra`).

Ideal for large-scale or monorepo Go applications, this tool helps ensure architectural integrity by preventing
undesired package-level coupling.

## Features

- 📊 Visualize internal package dependencies (Mermaid, Graphviz, HTML)
- 🚨 Detect and report invalid imports based on custom rules
- ✅ Output violations with actionable messages
- 🔍 Highlight architectural drift in pull requests
- 🧠 Perfect for layered, hexagonal, or clean architecture

## Installation

```bash
# Install goimportmaps into your project
go get -tool github.com/mickamy/goimportmaps@latest

# or install it globally
go install github.com/mickamy/goimportmaps@latest
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

| Option                        | Description                                             |
|-------------------------------|---------------------------------------------------------|
| `--format=mermaid`            | Output format: `text`, `mermaid`, `html`, or `graphviz` |

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

If `handler` imports `infra` directly, the tool will detect:

```bash
🚨 Violation detected: handler imports infra
🛠️  Fix: handler should depend on service, not infra directly
```

### Mermaid Output

```
```mermaid
graph TD
  main --> handler
  handler --> service
  service --> infra
  handler --> infra %% ❌
```

```

## Configuration

`.goimportmaps.yaml`

```yaml
forbidden:
  - from: handler
    to: infra
  - from: app
    to: db
```

## HTML Output

Use `--format=html` to generate a standalone static report:

```bash
goimportmaps ./... --format=html > report.html
```

This report can be viewed in your browser or uploaded as a CI artifact.

## License

[MIT](./LICENSE)
