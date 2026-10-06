# NFS

> [→ Enumeration Techniques](CHEATSHEETS/ENUMERATION/SERVICES/nfs.md)

## Overview

NFS (Network File System) is a distributed file system protocol that allows a client to access files over a network as if they were local. It is native to Unix/Linux environments and is commonly used for shared storage in enterprise Linux infrastructure and HPC clusters.

## Typical Targets

- Linux/Unix servers and storage appliances (NetApp, EMC, ZFS-based)
- High-performance computing clusters
- Development and build servers with shared home directories
- Internal file servers in mixed Linux environments

## Common Ports

| Port | Protocol | Usage |
|------|----------|-------|
| 111 | TCP/UDP | rpcbind / portmapper (required for NFS) |
| 2049 | TCP/UDP | NFS main port |

## Common Attack Paths

- **Unauthenticated Access to Exports:** NFS exports configured with `*(rw)` or broad network ranges allow any host on that range to mount the share without any authentication.
- **no_root_squash Misconfiguration:** By default NFS maps the remote root user to `nobody`. When `no_root_squash` is set, a root user on the client machine is trusted as root on the share — allowing an attacker to plant SUID binaries or modify any file.
- **Sensitive File Exposure:** NFS shares often contain SSH keys, configuration files, application secrets, or home directories with `.bash_history` and credential files.
- **Pivot via Planted Files:** Writing an SSH public key into a mounted `/home/user/.ssh/authorized_keys` allows passwordless SSH access to the server as that user.
