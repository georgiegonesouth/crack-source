# CLAUDE.md

This file provides guidance to Claude Code when working with this repository.

## Project overview

**crack-source** is a dual-interface penetration testing reference tool. It presents a structured knowledge base of attack techniques, enumeration cheatsheets, and tool references through two complementary interfaces:

- **`crack-source-web/`** — browser-based UI, zero dependencies, served statically
- **`crack-source/`** — terminal UI written in Go, runs natively in the operator's shell

Both interfaces consume the same content from `CONTENT/KNOWLEDGE/` and the same manifest from `CONTENT/MANIFEST/manifest.json`. They are peers, not one being a subset of the other.

The TUI is the primary development focus — it is being built toward full feature parity with the web UI plus terminal-native extras (shell drop, embedded terminal, session log, etc.).

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

Used inside `.md` files — parsed by both the web UI and the TUI's token system:

- `<!-- token-options:TOKEN\noption1\noption2 -->` — render "Use" buttons for predefined values
- `<!-- token-table:TOKEN -->` — table cells become "Use" buttons
- `<!-- token-section:TOKEN -->` — code block options
- `<!-- token-exclude:TOKEN1,TOKEN2 -->` — exclude tokens from the local input bar

---

## Web UI (`crack-source-web/`)

### Running locally

```bash
# From repo root:
python3 -m http.server 8080
# Open: http://localhost:8080/crack-source-web/
```

The server must be started from the **repo root**, not from `crack-source-web/`. The app fetches `web.config` to discover the manifest path, then resolves content paths relative to the repo root via `REPO_ROOT=..`.

Direct `file://` open fails — browsers block `fetch()` for local files.

### Configuration (`crack-source-web/web.config`)

```
MANIFEST_PATH=CONTENT/MANIFEST/manifest.json
REPO_ROOT=..
```

`app.js` fetches this file at startup and uses the values to construct all content URLs. Paths are relative: `MANIFEST_PATH` is relative to the repo root; `REPO_ROOT` is relative to the `crack-source-web/` directory (i.e., `..` points up to the repo root).

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

## TUI (`crack-source/`)

### Building and running

```bash
cd elluminos-tui && go build && ./crack-source
```

### Configuration (`crack-source/tui.config`)

```
MANIFEST_PATH=/absolute/path/to/CONTENT/MANIFEST/manifest.json
REPO_ROOT=/absolute/path/to/repo
```

The TUI reads files directly from the filesystem (no HTTP), so absolute paths are required here. `REPO_ROOT` is joined with manifest entry paths to locate content files.

### Architecture

Built with Charmbracelet (bubbletea, glamour, lipgloss). Two sub-packages under `internal/`:

- `internal/manifest/manifest.go` — manifest loading; exports `Entry`, `LoadFiles`
- `internal/parser/blocks.go` — exports `Block`, `BlockKind`, `ParsedDoc`, `ParseDoc`, constants
- `internal/parser/cleanup.go` — exports `ReToken`, `StripFrontMatter`, `StripDirectives`, `StripANSI`, `ExtractLocalTokens`, `ExtractTokenOptionValue`, `InsertCursorMark`, `PrevWord`, `NextWord`

Root `package main` file layout:

| File | Responsibility |
|---|---|
| `main.go` | Entry point, config loading; `tea.WithMouseCellMotion()` enabled |
| `model_core.go` | `Model` struct, `newModel`, `Init` |
| `model_types.go` | All types and constants (`pane`, `appMode`, `tokenKeys`, layout consts) |
| `model_helpers.go` | Computed properties on `Model` (`bodyH`, `contentW`, etc.) |
| `update.go` | Main key/message handler (`Update`); `tea.MouseMsg` for cmd pane scroll |
| `view.go` | Notebook-mode render functions |
| `render.go` | Block renderer, lipgloss styles, `splitTips`, `headings` |
| `content.go` | File loading, token substitution, `runCmd`, `isCdCmd`/`resolveCd` |
| `search.go` | Sidebar search, content index, `filterEntries` |
| `profile.go` | Token profile persistence (`~/.crack-source/profiles/`) |
| `editor.go` | `editorState`, `dirEntry`, all buffer and dir operations |
| `editor_view.go` | Editor sidebar and content pane rendering |
| `editor_update.go` | Editor key handlers: `handleModeToggle`, `handleEditorKey`, `handleEditorSidebarKey`, `handleEditorContentKey`, `editorScrollUpdate` |
| `cmd.go` | Cmd pane shell: `cmdPromptPrefix`, `cmdUpdateViewport`, `cmdClear`, `handleCmdKey` |

### App modes

Toggle with `ctrl+e` from any pane (including paneCmd).

**Notebook mode** (`modeNotebook`, default) — manifest-based reference notebook:
- Collapsible sidebar nav; `k`/`j` scroll, `K`/`J` jump folder; `'` opens full-text search (two-phase: type to filter, `enter` → nav mode; `'` re-focuses input from nav mode; any printable key exits nav mode)
- Content pane: collapsible H2/H3, heading nav `←`/`→`; code block nav `d`/`f`, `esc` deselects, `x` drops block to cmd pane, `u` applies token-section value; manual edit `enter` (`alt+←/→`, `ctrl+a/e`; `enter` exits, `esc` discards); `tab`/`shift+tab` cycle local tokens within selected block
- Token bar: global tokens `TARGET_IP`, `PORT`, `LHOST`, `LPORT`, `USER`, `PASSWORD`, `DOMAIN`; page-local tokens per file
- Tips panel auto-extracted from `## Tips` section of each file

**Editor mode** (`modeEditor`) — filesystem file manager + text editor:
- Sidebar: `k`/`j` or arrows navigate; `enter`/`l` opens dir or file; `tab`/`right` jumps to editor pane when file is open; `ctrl+e` toggles back to notebook
- Content pane: arrow keys move cursor; `home`/`ctrl+a` line start, `end` line end; printable runes insert; `backspace`/`delete` delete; `enter` splits line; `tab` inserts 4 spaces; `ctrl+s` save; `esc` → sidebar; line numbers + block cursor; horizontal scroll when line exceeds pane width
- Token bar hidden in this mode

### Persistent working directory (`workDir`)

`workDir string` on `Model` is the shared working directory for both modes:
- All shell commands run with `cmd.Dir = workDir`
- `cd` in the cmd pane is intercepted — updates `workDir` (and reloads editor sidebar if in editor mode)
- Navigating into a directory in the editor sidebar updates `workDir`

### Cmd pane (both modes)

Terminal-like shell panel implemented in `cmd.go`. Key Model fields:
- `cmdScrollback string` — accumulated history (all past prompts + outputs)
- `cmdCurrentLine string` / `cmdCursorPos int` — current input with rune-level editing
- `cmdVp viewport.Model` — scrollable viewport; content = `cmdScrollback + prompt + cursor`; prompt format is `user@host dir %` (green user, cyan dir, white %)
- `cmdPromptAtTop bool` — set by `ctrl+l`/`clear`; appends `Height-1` trailing `\n` so `GotoBottom()` places prompt at top
- `cmdClearLine int` — scrollback line count at last clear; `YOffset` is floored here to hide pre-clear history

`cmdUpdateViewport()` always calls `GotoBottom()` then enforces `YOffset >= cmdClearLine`. Entering a command clears `cmdPromptAtTop`, returning to normal bottom-anchored flow.

**Key bindings:**
- `C` — focus from any pane; `esc` — return to content pane
- `enter` — run command; `cd` intercepted; `clear`/`reset`/`ctrl+l` — visual clear
- `ctrl+shift+c`/`ctrl+shift+v` — clipboard copy/paste
- `up`/`down` — history cycle; `alt+←/→`, `ctrl+w` — word navigation/delete
- `shift+up`/`shift+down` — resize; `ctrl+shift+up`/`ctrl+shift+down` — snap max/min
- Mouse wheel scrolls scrollback when paneCmd is focused

---

## Shared design constraints

- All content paths in manifests must be relative to the repo root — no absolute filesystem paths
- Both UIs use the same manifest and the same content files; format decisions in `.md` files affect both
- Token placeholder names are shared across both UIs; adding a new global token requires changes in both `app.js` (TOKENS array) and `main.go` (tokenKeys array)
- The `LIBRARY` subtree is the only part of the nav that is treated differently between the two UIs (separate panel in web, flat inclusion in TUI sidebar)
