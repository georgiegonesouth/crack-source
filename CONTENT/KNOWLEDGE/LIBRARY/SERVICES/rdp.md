# RDP

> [→ Enumeration Techniques](CHEATSHEETS/ENUMERATION/SERVICES/rdp.md)

## Overview

RDP (Remote Desktop Protocol) is Microsoft's proprietary protocol for graphical remote access to Windows systems. It provides a full desktop experience and is one of the most commonly exposed services on internet-facing Windows hosts. RDP is used by administrators, support staff, and legitimate remote workers.

## Typical Targets

- Internet-facing Windows servers and workstations
- Jump servers and bastion hosts
- Domain controllers and terminal servers
- Any Windows system with remote management enabled

## Common Ports

| Port | Protocol | Usage |
|------|----------|-------|
| 3389 | TCP | Standard RDP |
| 3389 | UDP | UDP transport (RDP 8.0+) |

## Common Attack Paths

- **Credential Brute Force / Password Spraying:** RDP is a prime target for brute-force attacks. Exposed RDP is continuously scanned and attacked by automated botnets.
- **BlueKeep (CVE-2019-0708):** Pre-auth RCE vulnerability in Windows 7 / Server 2008. Widely patched but still found on unpatched legacy systems.
- **DejaBlue (CVE-2019-1181/1182):** Similar pre-auth RCE vulnerabilities affecting Windows 8.1 and Server 2012.
- **Pass-the-Hash with Restricted Admin Mode:** When Restricted Admin Mode is enabled, it is possible to authenticate to RDP using an NTLM hash instead of a plaintext password.
- **MitM with RDP Proxy Tools:** Tools like Seth or pyrdp can intercept RDP sessions to capture credentials or inject keystrokes when a victim connects through a rogue gateway.
- **Session Hijacking:** On a compromised host, administrators can attach to active or disconnected RDP sessions of other users without credentials using `tscon`.
