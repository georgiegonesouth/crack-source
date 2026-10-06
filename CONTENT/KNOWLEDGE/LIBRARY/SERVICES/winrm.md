# WinRM

> [→ Enumeration Techniques](CHEATSHEETS/ENUMERATION/SERVICES/winrm.md)

## Overview

WinRM (Windows Remote Management) is Microsoft's implementation of the WS-Management protocol, enabling remote PowerShell sessions and command execution on Windows systems. It is enabled by default on Windows Server 2012+ and is the backbone of PowerShell Remoting (`Enter-PSSession`, `Invoke-Command`). It is a primary lateral movement path in Windows environments.

## Typical Targets

- Windows Server 2012 and later (WinRM enabled by default)
- Domain controllers and member servers in Active Directory environments
- Azure VMs and cloud-hosted Windows instances
- Jump servers in managed Windows environments

## Common Ports

| Port | Protocol | Usage |
|------|----------|-------|
| 5985 | TCP | WinRM over HTTP |
| 5986 | TCP | WinRM over HTTPS |

## Common Attack Paths

- **Credential Brute Force:** WinRM can be brute-forced with tools like CrackMapExec. The service has no built-in lockout independent of the OS account lockout policy.
- **Pass-the-Hash with evil-winrm:** NTLM hash authentication is supported; a captured hash from any technique (LSASS dump, NTLM relay, secretsdump) can be used directly without cracking.
- **Lateral Movement with Compromised Credentials:** WinRM is the go-to method for interactive lateral movement once valid credentials or hashes are obtained.
- **Kerberos Ticket (Pass-the-Ticket):** Valid Kerberos tickets (e.g., from a golden/silver ticket attack or TGT/TGS theft) can authenticate WinRM sessions.
- **Misconfigured HTTP Listener:** WinRM over plain HTTP (port 5985) transmits session data unencrypted; on a network where traffic can be captured, this exposes commands and output.
