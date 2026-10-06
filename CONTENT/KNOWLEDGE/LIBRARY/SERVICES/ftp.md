# FTP

> [→ Enumeration Techniques](CHEATSHEETS/ENUMERATION/SERVICES/ftp.md)

## Overview

FTP (File Transfer Protocol) is a legacy protocol for transferring files between a client and server. It transmits credentials and data in clear text, making it inherently insecure. Despite this, it remains common on older systems, embedded devices, and some internal file-sharing setups.

## Typical Targets

- Legacy servers and internal file shares
- NAS devices and embedded systems
- Development or staging servers
- Industrial control systems and older network appliances

## Common Ports

| Port | Protocol | Usage |
|------|----------|-------|
| 21 | TCP | Command/control channel |
| 20 | TCP | Active mode data transfer |
| Ephemeral | TCP | Passive mode data transfer (PASV) |

## Common Attack Paths

- **Anonymous Login:** Many FTP servers are misconfigured to allow unauthenticated anonymous access, which may expose sensitive files or allow file uploads.
- **Credential Brute Force:** FTP has no built-in lockout by default, making it susceptible to password spraying and brute-force attacks.
- **Clear-text Credential Sniffing:** FTP sends usernames and passwords in plain text; anyone on the network path can capture them.
- **Misconfigured Write Permissions:** If anonymous or authenticated users can write files, an attacker can upload webshells, backdoors, or replace existing binaries.
- **FTP Bounce Attack:** Using the PORT command to proxy connections through the FTP server to scan or connect to other internal hosts.
- **Known CVEs:** Older FTP daemons (vsftpd 2.3.4 backdoor, ProFTPD vulnerabilities) have well-known exploits that grant remote code execution.
