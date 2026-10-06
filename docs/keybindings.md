# Keybindings (TUI)

All keys for the terminal UI (`crack-source-tui/`). Keys are grouped by where they apply.
Any of these can be remapped — see [configuration.md](configuration.md#custom-keybindings).

> The canonical source for every binding is `crack-source-tui/keybindings.yaml`. If this
> doc and that file ever disagree, the file wins.

## Application

| Action | Default key |
|---|---|
| Quit | `ctrl+c`, `q` |
| Toggle notebook ↔ editor mode | `ctrl+e` |

## Pane focus

| Action | Default key |
|---|---|
| Focus cmd pane | `C` |
| Focus token bar | `t` |
| Focus search | `'` (also re-focuses search input from nav mode) |
| Focus sidebar | `h` |
| Focus content pane | `tab`, `l` |

## Notebook — sidebar

| Action | Default key |
|---|---|
| Move down | `k`, `down` |
| Move up | `j`, `up` |
| Jump to next folder | `K` |
| Jump to previous folder | `J` |

## Notebook — content pane

| Action | Default key |
|---|---|
| Next section (heading) | `right`, `k` |
| Previous section (heading) | `left`, `j` |
| Scroll up | `l`, `up` |
| Scroll down | `ö`, `down` |
| Next code block | `f` |
| Previous code block | `d` |
| Drop selected code block to cmd pane | `x` |
| Apply token-section value | `u` |

## Cmd pane

| Action | Default key |
|---|---|
| Run command | `enter` |
| Clear | `ctrl+l` |
| History up / down | `up` / `down` |
| Copy current line | `ctrl+shift+c` |
| Paste | `ctrl+shift+v` |

### Cmd pane height

| Action | Default key |
|---|---|
| Grow / shrink | `shift+up` / `shift+down` |
| Snap to max / min | `ctrl+shift+up` / `ctrl+shift+down` |

## Editor mode

| Action | Default key |
|---|---|
| Sidebar up / down | `j` / `k` (or arrows) |
| Open dir or file | `enter`, `l` |
| Save | `ctrl+s` |

> Note: in editor mode the sidebar `j`/`k` directions are inverted relative to the
> notebook sidebar.

## Shared editing primitives

Used across the cmd pane, inline code-edit mode, and the editor content pane:

| Action | Default key |
|---|---|
| Cancel / escape | `esc` |
| Confirm | `enter` |
| Indent / outdent | `tab` / `shift+tab` |
| Backspace | `backspace`, `ctrl+h` |
| Delete | `delete` |
| Move cursor | arrow keys |
| Line start / end | `home` (`ctrl+a`) / `end` |
| Word left / right | `alt+left` / `alt+right` |
| Delete word | `ctrl+w` |
