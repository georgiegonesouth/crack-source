# Configuration

Both interfaces are configured by a small plain-text config file, plus an optional
keybindings file for the TUI.

---

## Web UI — `crack-source-web/web.config`

```
MANIFEST_PATH=CONTENT/MANIFEST/manifest.json
REPO_ROOT=..
```

- `MANIFEST_PATH` — path to the manifest, relative to the repo root.
- `REPO_ROOT` — path to the repo root, relative to the `crack-source-web/` directory
  (`..` points one level up).

`app.js` fetches this at startup and uses it to build all content URLs. Because of this,
the static server **must be started from the repo root**:

```bash
python3 -m http.server 8080
# open http://localhost:8080/crack-source-web/
```

---

## TUI — `crack-source-tui/tui.config`

```
MANIFEST_PATH=../CONTENT/MANIFEST/manifest.json
REPO_ROOT=..
```

- `MANIFEST_PATH` — path to the manifest.
- `REPO_ROOT` — joined with each manifest entry path to locate content files.

The TUI reads files directly from the filesystem (no HTTP). The shipped values are
relative to where the binary runs (i.e. from inside `crack-source-tui/`). Absolute paths
also work if you prefer to run the binary from anywhere.

<!-- TODO: confirm/adjust if you standardize on absolute paths -->

---

## Custom keybindings (TUI)

Every TUI key can be remapped. The defaults live in `crack-source-tui/keybindings.yaml`;
to override them, copy that file to:

```
~/.config/crack-source/keybindings.yaml
```

Each entry maps an action to a list of keys. **Any action you omit keeps its built-in
default**, so your override file only needs the bindings you actually want to change.

```yaml
# example: use vim-style quit and a different mode-toggle
quit: ["ctrl+c", "ZZ"]
mode_toggle: ["ctrl+space"]
```

See [keybindings.md](keybindings.md) for the full list of action names and their defaults.

---

## Adding a new global token

A global token must be registered in **both** UIs:

- Web: the `TOKENS` array in `crack-source-web/app.js`
- TUI: the `tokenKeys` array in `crack-source-tui/main.go`

Token placeholder names are shared across both interfaces, so keep them identical.
