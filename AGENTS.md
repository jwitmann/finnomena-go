# AGENTS.md

Guidelines and context for AI agents working on the codebase.

## LLM Guidelines

### Do

- Run `go build` after any code change to verify it compiles
- Run `go test ./...` before committing to ensure tests pass
- Keep existing functionality unchanged by default
- **Parallelize operations where possible** — use goroutines, worker pools, batch processing
- **Prefer caching** — cache computed results, API responses, and expensive operations; use TTL-based invalidation
- **Write deterministic concurrent code** — when aggregating results from goroutines or iterating over maps, use sorted keys or a stable ordering so repeated runs produce identical outputs (e.g., map iteration order and channel completion order are non-deterministic in Go)
- Commit changes locally after task completion

### Don't

- Modify scoring/analysis algorithms without explicit user request
- Remove existing features or API endpoints
- Change external API behavior (Finnomena.com)
- Break backward compatibility in settings/storage
- Push commits to remote — commit locally only, never push

### Ask First

- Refactoring that changes behavior
- Adding new dependencies
- Removing functionality
- Breaking changes to data formats

### Stay on Goal — Park Side Findings, Don't Drift

While working on a task you will notice things unrelated to the goal: dead code, stale docs, schema/code mismatches, suspicious patterns, potential bugs, cleanup opportunities, or a nice-to-have optimization you deliberately skipped to keep the change minimal (e.g., "add index on `group_name` for distinct query"). **Do not act on them. Do not expand scope. Do not stop to ask the user which option they prefer unless the current task is blocked on the answer.**

Instead, **park** the finding so it is not lost and continue with the assigned task:

- Append a one-line entry to `docs/NOTES.md` (create the file if missing) under a dated `## YYYY-MM-DD` heading.
- Format: `- [location] short description (found while: <task>)`. Examples:
  - `- [internal/db/queries.go:ArchiveOldTorrents] dead code, never called from main/aggregator/server (found while: updating docs)`
  - `- [internal/db/schema.sql] consider adding index on group_name for DISTINCT queries (found while: <task>)`
- One line per finding. No analysis, no proposed fixes, no questions in the file — those go in a later, dedicated task.
- Resume the original task immediately after writing the note. Do not switch files, run exploratory commands, or open related code "just to check".
- **At the end of the task**, list any notes you parked during this run in your final summary to the user (just the one-line entries — no analysis or proposed actions). If nothing was parked, say so explicitly. This surfaces findings without acting on them.
- **When a parked finding is resolved, delete its entry from `docs/NOTES.md`.** Do not leave `FIXED:` annotations or accumulate historical entries; the file is a parking lot for open findings, not a dump of past issues.

Rationale: context is finite and drift compounds. A task to "update docs" should never turn into "should I delete the archive system?" mid-flight. Noteworthy observations are valuable, but they belong in a parking lot, not in the critical path of the current task. The user can triage `docs/NOTES.md` whenever they choose and spin off focused tasks from it.

## Browse Tool — Agent Reference

### Overview

The `scripts/browse` tool is a lightweight wrapper around Chromium that fetches webpage content while bypassing bot detection/WAF protection.

### When to Use

Use `browse` instead of `curl`/`wget` when:

- Site returns "Access Denied" or "Blocked"
- Site uses Cloudflare, Akamai, or similar WAF
- Site is a Single Page Application (SPA) that requires JavaScript rendering
- Standard HTTP tools return empty or bot challenge pages

### Usage

```bash
# Basic usage
browse "https://example.com"

# Save to file
browse "https://example.com" > output.html

# Pipe to other tools
browse "https://example.com" | grep "pattern"

# For JSON APIs
browse "https://api.example.com/data" | jq '.'
```

## Critical Dev Rules

- Pure Go only — no CGO
- **Binaries ONLY in `bin/`** — never build to the repo root. The Makefile builds to `bin/`; do not create root-level binaries
- **`.gitignore` patterns must be anchored** — use `/pattern` (root-relative) for top-level files/dirs to avoid accidentally matching subpackages. Example: use `/db/` not `db/` (the latter matches `internal/db/` too)

## Tech Stack

- **Language:** Go 1.25 (pure Go, no CGO/MinGW required)
- **Data Storage:** Feather format (Apache Arrow) for price data, JSON for config

## Developer Tools

The following tools are installed in `~/go/bin/` and require `export PATH="$HOME/go/bin:$PATH"` to run:

| Tool | Purpose | Run |
|------|---------|-----|
| `air` | Live-reload development | `air` (uses `.air.toml`) |
| `golangci-lint` | Comprehensive linter | `golangci-lint run ./...` |
| `staticcheck` | Static analysis | `staticcheck ./...` |
| `gofumpt` | Stricter `gofmt` | `gofumpt -w .` |
| `gocyclo` | Cyclomatic complexity | `gocyclo -over 15 ./...` |
| `dupl` | Duplicate code detection | `dupl -t 100 ./...` |
| `dlv` | Delve debugger | `dlv debug ./cmd/thaifa` |
| `govulncheck` | Vulnerability scanner | `govulncheck ./...` |
| `gotests` | Test boilerplate generator | `gotests -all -w file.go` |
| `richgo` | Colorized test output | `richgo test ./...` |

### Live-reload with Air

```bash
export PATH="$HOME/go/bin:$PATH"
air
```

This watches `.go` and `.html` files and rebuilds via `make thaifa BIN_DIR=tmp`, preserving version LDFLAGS.

## Code Conventions

### Go Style
- Standard Go formatting (gofmt)
- No comments unless explicitly requested
- Error handling with fmt.Errorf and %w for wrapping
- Mutex-protected settings with RLock/Lock
