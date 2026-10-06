# Authoring content

All content is plain Markdown under `CONTENT/KNOWLEDGE/`. Both UIs read the same files, so
anything you write here shows up identically in the web app and the TUI.

After **adding, renaming, or removing** any file, regenerate the manifest:

```bash
cd CONTENT/MANIFEST && python3 generate_manifest.py
```

Manifest paths are always relative to the repo root — never absolute.

---

## Placeholder tokens

Tokens are angle-bracket placeholders that the UI substitutes with values the operator
fills in once. Unfilled tokens are highlighted so you can spot them.

### Global tokens

Filled in the header/token bar and reused across every page:

| Token | Meaning |
|---|---|
| `<TARGET_IP>` | Target machine IP or FQDN |
| `<PORT>` | Target service port |
| `<LHOST>` | Listener / VPN IP |
| `<LPORT>` | Listener port |
| `<USER>` | Username |
| `<PASSWORD>` | Password |
| `<DOMAIN>` | Domain |

> Use these exact names. Variants like `<IP>`, `<TARGET>`, or `<FQDN/IP>` are **not**
> substituted.

### Page-local tokens

Any other UPPERCASE angle-bracket placeholder is treated as page-local and gets its own
input on that page's token bar (e.g. `<WORDLIST>`, `<RULE>`, `<SOURCE>`).

Page-local tokens that refer to small inline values use lowercase, e.g. `<database>`,
`<share>`, `<user>`.

---

## Markdown directives

Directives are HTML comments the UIs parse to add interactivity. They have no effect when
the file is viewed as plain Markdown (e.g. on GitHub), so they're safe to leave in.

| Directive | Effect |
|---|---|
| `<!-- token-options:TOKEN\noption1\noption2 -->` | Render "Use" buttons for predefined values |
| `<!-- token-table:TOKEN -->` | Turn table cells into "Use" buttons |
| `<!-- token-section:TOKEN -->` | Offer code-block values for a token |
| `<!-- token-exclude:TOKEN1,TOKEN2 -->` | Exclude tokens from the local input bar |

<!-- TODO: add a worked example of each directive if you want fuller docs -->

---

## Tips section

A `## Tips` heading in a file is auto-extracted and shown in a dedicated panel (right
sidebar in the web UI, Tips panel in the TUI). Put page-specific guidance there.

---

## The LIBRARY subtree

`CONTENT/KNOWLEDGE/LIBRARY/` is treated specially: the web UI gives it a separate
Library panel, while the TUI folds it into the main sidebar. PDFs, writeups, and external
references live here.
