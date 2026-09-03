# RSYNC

> [→ Enumeration Techniques](CHEATSHEETS/ENUMERATION/SERVICES/rsync.md)

## Overview

rsync is a fast, incremental file transfer utility widely used for backups, mirroring, and deployment on Unix/Linux systems. When run as a daemon (rsyncd), it listens on a TCP port and exposes named "modules" (directories) that clients can sync to or from. It is frequently found on internal servers and backup infrastructure.

## Typical Targets

- Linux/Unix backup and storage servers
- CI/CD and deployment servers
- Internal file sync servers

## Common Ports

| Port | Protocol | Usage |
|------|----------|-------|
| 873 | TCP | rsync daemon |

## Common Attack Paths

- **Unauthenticated Module Access:** rsync modules are often configured without a password, allowing any client to list and download their entire contents — including config files, SSH keys, database dumps, and source code.
- **Credential Brute Force:** When authentication is enabled, rsync uses a simple challenge-response mechanism with no lockout, making it brute-forceable.
- **Arbitrary File Read:** Downloading a module with sensitive files (`/home`, `/etc`, `/var/backups`) can yield private SSH keys, shadow files, or application credentials.
- **Arbitrary File Write:** If a module is writable, an attacker can upload files — planting an SSH `authorized_keys`, cron jobs, or replacing binaries to achieve persistence or code execution.
