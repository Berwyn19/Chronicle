# Chronicle — Design

Chronicle is a local-first CLI that records AI-assisted development history. It
does **not** replace Git: Git remains the source of truth for source code.
Chronicle records the *intent* behind each AI interaction (prompt, model,
metadata) and stores a Git patch of the resulting changes.

## Philosophy

Follows the spirit of Git, ripgrep, tmux, and Lazygit: small, composable, zero
configuration, local-first, fast, Unix-y.

## Technology

- **Go** (1.24)
- **Cobra** for the CLI
- **SQLite** for metadata, via `modernc.org/sqlite` (pure Go, no cgo)
- Shell out to the **Git CLI** — never reimplement Git
- No web server, no cloud, no telemetry

### Why `modernc.org/sqlite`

Pure-Go driver, so the binary is self-contained and cross-compiles trivially
with no C toolchain. The common cgo alternative (`mattn/go-sqlite3`) would
require a C compiler at build time, working against the "small / zero-config"
goal.

## On-disk layout

Everything lives under a single `.chronicle/` directory in the project root:

```
.chronicle/
    metadata.db          SQLite database (intent + pointers)
    events/
        <event-id>.patch  one Git patch per AI interaction
```

Patches are stored as plain `git diff` output. File snapshots are **never**
stored — Git already has the content.

## SQLite schema

```
events
--------
id          TEXT PRIMARY KEY
prompt      TEXT NOT NULL
model       TEXT NOT NULL
timestamp   TEXT NOT NULL
patch_path  TEXT NOT NULL
commit_hash TEXT
```

Scope is per-project: each project has its own `.chronicle/metadata.db`, so the
`events` table holds exactly that project's history (the same way `.git/`
scopes a Git repo by location). There is an index on `events(timestamp)` for
listing in chronological order.

## Project structure

```
main.go                      thin entrypoint -> cmd.Execute()
cmd/
    root.go                  Cobra root command; no business logic
    init.go                  chronicle init
internal/
    paths/paths.go           owns the .chronicle/ layout
    store/store.go           owns SQLite; the only DB-aware package
    store/schema.sql         embedded schema (//go:embed)
```

### Package responsibilities

Each package has a single responsibility:

- **`cmd`** — Cobra wiring only. Each command is a thin wrapper that delegates
  to the internal packages. Adding a command means adding one file here.
- **`internal/paths`** — the single source of truth for the directory layout.
  Nothing else hard-codes `metadata.db`, `events/`, or `.patch`.
- **`internal/store`** — the only package that talks to SQLite. The schema is
  embedded and applied idempotently on `Open`, so "open the DB" and
  "initialize the DB" are the same code path. No separate migration machinery
  (avoiding premature abstraction).

## Commands (v0)

Exactly five commands. No more.

| Command                    | Status      | Purpose                                   |
| -------------------------- | ----------- | ----------------------------------------- |
| `chronicle init`           | implemented | Create `.chronicle/` and the database     |
| `chronicle list`           | planned     | List recorded events                      |
| `chronicle show <id>`      | planned     | Show metadata for one event               |
| `chronicle diff <id>`      | planned     | Show the patch for one event              |
| `chronicle file <path>`    | planned     | Show events that touched a file           |

### Explicitly out of scope

replay, cloud sync, GitHub integration, VS Code extension, semantic diff,
prompt blame, dashboards, analytics, embeddings, vector databases.

## `chronicle init` behavior

1. Resolve `.chronicle/` relative to the current working directory.
2. If it already exists, print a message and exit successfully (safe re-run).
3. Otherwise create `.chronicle/` and `.chronicle/events/`.
4. Open the SQLite database, which applies the schema.

## Open design question

The schema and patch flow imply something must *write* sessions and events
(`git diff` -> `events/<id>.patch` -> a row). The v0 command list is otherwise
all read-side. How events get recorded (a dedicated path vs. external
integration) is still to be decided.
