# SSH

> [→ Enumeration Techniques](CHEATSHEETS/ENUMERATION/SERVICES/ssh.md)

## Overview

SSH (Secure Shell) is the standard encrypted protocol for remote command-line access to Unix/Linux systems and network devices. It supports password authentication, public-key authentication, and port forwarding. SSH is found on virtually every Linux server and most network infrastructure.

## Typical Targets

- Linux and Unix servers (essentially all of them)
- macOS systems
- Network devices (routers, switches, firewalls)
- Embedded systems (IoT, NAS, access points)

## Common Ports

| Port | Protocol | Usage |
|------|----------|-------|
| 22 | TCP | Standard SSH |

## Common Attack Paths

- **Credential Brute Force:** Despite being encrypted, SSH login attempts can be automated with tools like Hydra or Medusa. Servers without fail2ban or rate limiting are vulnerable.
- **Password Spraying:** Testing a small set of common passwords against many accounts avoids lockouts while still finding weak credentials.
- **Weak or Reused Private Keys:** Private SSH keys found in `.bash_history`, on SMB shares, in Git repos, or in web application directories may be valid on the target host.
- **Authorized Keys Injection:** With write access to a user's `~/.ssh/authorized_keys` (via NFS, FTP, rsync, or file write vulnerability), an attacker can add their own public key for persistent access.
- **Username Enumeration (Older OpenSSH):** OpenSSH versions before 7.7 responded differently to valid vs. invalid usernames during key-based authentication.
- **SSH Tunneling for Pivoting:** SSH's `-L`, `-R`, and `-D` forwarding options are widely used to pivot into internal networks through a compromised host.
- **Known CVEs:** Older OpenSSH versions have had critical vulnerabilities (e.g., CVE-2023-38408 ssh-agent RCE, CVE-2024-6387 "regreSSHion").
