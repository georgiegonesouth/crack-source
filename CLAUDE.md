# CLAUDE.md

This file provides guidance to Claude Code when working with this repository.

## Project overview

**Elluminos** is a dual-interface penetration testing reference tool. It presents a structured knowledge base of attack techniques, enumeration cheatsheets, and tool references through two complementary interfaces:

- **`elluminos-web/`** — browser-based UI, zero dependencies, served statically
- **`elluminos-tui/`** — terminal UI written in Go, runs natively in the operator's shell

Both interfaces consume the same content from `CONTENT/KNOWLEDGE/` and the same manifest from `CONTENT/MANIFEST/manifest.json`. They are peers, not one being a subset of the other.

## Content structure

```
CONTENT/
├── KNOWLEDGE/          ← all cheatsheet and reference content
│   ├── ACTIVE_DIRECTORY/
│   ├── ENUMERATION/
│   ├── FILE_TRANSFER/
│   ├── LIBRARY/        ← PDFs, writeups, external references
│   ├── OS/
│   ├── SHELLS/
│   └── MANIFEST/       ← generate_manifest.py for this subtree
└── MANIFEST/           ← canonical manifest used by both UIs
    ├── generate_manifest.py
    └── manifest.json   ← auto-generated, do not edit by hand
```

### Regenerating the manifest

Run from `CONTENT/MANIFEST/` after adding, renaming, or removing any file under `CONTENT/KNOWLEDGE/`:

```bash
cd CONTENT/MANIFEST && python3 generate_manifest.py
```

Manifest paths are **relative to the repo root** (e.g. `CONTENT/KNOWLEDGE/ACTIVE_DIRECTORY/ATTACKS/adcs_attacks.md`). Do not use absolute paths; the generate script computes the scan root from its own location.

### Placeholder convention for cheatsheet files

| Token | Meaning |
|---|---|
| `<TARGET_IP>` | Target machine IP or FQDN |
| `<PORT>` | Target service port |
| `<LHOST>` | Listener / VPN IP |
| `<LPORT>` | Listener port |

Do not use `<IP>`, `<TARGET>`, `<FQDN/IP>`, or other variants — they won't be substituted.

Page-local tokens use lowercase: `<database>`, `<share>`, `<user>`, etc.

### Markdown comment directives

Used inside `.md` files to enhance the web UI's token substitution system:

- `<!-- token-options:TOKEN\noption1\noption2 -->` — render "Use" buttons for predefined values
- `<!-- token-table:TOKEN -->` — table cells become "Use" buttons
- `<!-- token-section:TOKEN -->` — code block options
- `<!-- token-exclude:TOKEN1,TOKEN2 -->` — exclude tokens from the local input bar

---

## Web UI (`elluminos-web/`)

### Running locally

```bash
# From repo root:
python3 -m http.server 8080
# Open: http://localhost:8080/elluminos-web/
```

The server must be started from the **repo root**, not from `elluminos-web/`. The app fetches `web.config` to discover the manifest path, then resolves content paths relative to the repo root via `REPO_ROOT=..`.

Direct `file://` open fails — browsers block `fetch()` for local files.

### Configuration (`elluminos-web/web.config`)

```
MANIFEST_PATH=CONTENT/MANIFEST/manifest.json
REPO_ROOT=..
```

`app.js` fetches this file at startup and uses the values to construct all content URLs. Paths are relative: `MANIFEST_PATH` is relative to the repo root; `REPO_ROOT` is relative to the `elluminos-web/` directory (i.e., `..` points up to the repo root).

### Architecture

No build step, no framework — plain HTML/CSS/JS.

- `index.html` — entry point
- `app.js` — all application logic
- `style.css` — dark theme
- `web.config` — runtime config (manifest path, repo root)

### Rendering pipeline (`app.js`)

1. Fetch and parse `web.config` → populate `CONFIG`
2. Fetch manifest → build sidebar nav tree (non-LIBRARY nodes) and library panel (LIBRARY subtree)
3. Build full-text content index by prefetching all `.md` files
4. On file select: fetch markdown → strip front matter → parse with marked.js → syntax-highlight with Highlight.js
5. Extract page-local tokens from HTML comments → render per-page input bar
6. Apply global + local token substitution to highlighted HTML
7. Extract Tips section → right sidebar; make H2/H3 collapsible; add copy/reset buttons to code blocks

### Token substitution system

Global tokens (stored in `localStorage`): `<TARGET_IP>`, `<PORT>`, `<LHOST>`, `<LPORT>` — filled via the header input bar and persisted across sessions.

Page-local tokens: uppercase angle-bracket placeholders not in the global set, extracted from the markdown source and rendered as a per-page input bar.

Unfilled tokens render with red highlighting. The Copy button copies fully-substituted plain text.

### Navigation

Two trees are built from the manifest:
- **Main sidebar** — everything except the `LIBRARY` subtree; searched with `/`
- **Library panel** — the `LIBRARY` subtree; toggled with the Library button or `Ctrl+/`

Both support full-text search across file names and file contents.

---

## TUI (`elluminos-tui/`)

### Current state

Built with Charmbracelet (bubbletea, glamour, lipgloss). Implements:

- Collapsible sidebar navigation from the same manifest
- Markdown rendering via glamour with token substitution
- Global token input bar (`TARGET_IP`, `PORT`, `LHOST`, `LPORT`)
- Inline command pane: type a shell command, capture and display its output within the TUI

### Building and running

```bash
cd elluminos-tui && go build && ./pt-tui
```

### Configuration (`elluminos-tui/tui.config`)

```
MANIFEST_PATH=/absolute/path/to/CONTENT/MANIFEST/manifest.json
REPO_ROOT=/absolute/path/to/repo
```

The TUI reads files directly from the filesystem (no HTTP), so absolute paths are required here. `REPO_ROOT` is joined with manifest entry paths to locate content files.

### Target feature set (parity + terminal-native extras)

The TUI is intended to reach full feature parity with the web UI and go beyond it with capabilities that only make sense in a terminal context:

**Parity with web UI:**
- Full sidebar nav with search/filter
- Token substitution (global + page-local)
- All markdown comment directives (`token-options`, `token-table`, `token-section`, `token-exclude`)
- Tips panel (rendered separately from main content)
- Collapsible H2/H3 sections

**Terminal-native features:**
- **Shell drop** — execute any code block or substituted command directly in the system shell (via pty), with output streaming back into the TUI; no manual copy-paste
- **Embedded terminal** — a full interactive shell pane within the TUI for running tools, with the ability to switch between the reference pane and the shell without leaving the application
- **Session log** — append timestamped commands and outputs to a per-session log file automatically
- **Note taking** — inline editor for attaching operator notes to any cheatsheet file (stored separately from content, never modifying source files)
- **Engagement log** — structured logging of findings, commands run, and timestamps exportable to markdown

---

## Shared design constraints

- All content paths in manifests must be relative to the repo root — no absolute filesystem paths
- Both UIs use the same manifest and the same content files; format decisions in `.md` files affect both
- Token placeholder names are shared across both UIs; adding a new global token requires changes in both `app.js` (TOKENS array) and `main.go` (tokenKeys array)
- The `LIBRARY` subtree is the only part of the nav that is treated differently between the two UIs (separate panel in web, flat inclusion in TUI sidebar)
