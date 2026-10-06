#!/usr/bin/env python3
"""
Scans CONTENT/KNOWLEDGE/ and writes manifest.json for the cheatsheet web app.
Run whenever you add, rename, or remove files:

    python3 generate_manifest.py

Paths in the manifest are relative to the repo root so the web app can fetch
them directly when the server is started from the repo root.
"""

import os
import json
from datetime import datetime

# Relative path prefix written into manifest entries (from repo root)
PATH_PREFIX = "CONTENT/KNOWLEDGE"
EXTENSIONS = {".md", ".pdf"}
SKIP = {".DS_Store", "__pycache__", ".git", ".gitignore"}


def scan_dir(abs_path, rel_path):
    entries = []
    try:
        items = sorted(os.listdir(abs_path), key=str.lower)
    except PermissionError:
        return entries

    dirs, files = [], []
    for item in items:
        if item in SKIP or item.startswith("."):
            continue
        item_abs = os.path.join(abs_path, item)
        item_rel = f"{rel_path}/{item}"
        if os.path.isdir(item_abs):
            dirs.append((item, item_abs, item_rel))
        elif os.path.isfile(item_abs):
            if os.path.splitext(item)[1].lower() in EXTENSIONS:
                files.append((item, item_abs, item_rel))

    for name, abs_p, rel_p in dirs:
        entries.append({
            "name": name,
            "type": "dir",
            "path": rel_p,
            "children": scan_dir(abs_p, rel_p),
        })

    for name, abs_p, rel_p in files:
        entries.append({
            "name": name,
            "type": "file",
            "path": rel_p,
        })

    return entries


def count_files(tree):
    n = 0
    for node in tree:
        if node["type"] == "file":
            n += 1
        else:
            n += count_files(node.get("children", []))
    return n


def main():
    script_dir = os.path.dirname(os.path.abspath(__file__))
    # KNOWLEDGE/ is a sibling of this script's directory (MANIFEST/)
    content_abs = os.path.normpath(os.path.join(script_dir, "..", "KNOWLEDGE"))

    if not os.path.isdir(content_abs):
        print(f"ERROR: {content_abs} not found.")
        raise SystemExit(1)

    tree = scan_dir(content_abs, PATH_PREFIX)
    manifest = {
        "generated": datetime.now().isoformat(timespec="seconds"),
        "root": PATH_PREFIX,
        "tree": tree,
    }

    output_path = os.path.join(script_dir, "manifest.json")
    with open(output_path, "w", encoding="utf-8") as f:
        json.dump(manifest, f, indent=2, ensure_ascii=False)

    n = count_files(tree)
    print(f"manifest.json written — {n} files indexed.")


if __name__ == "__main__":
    main()
