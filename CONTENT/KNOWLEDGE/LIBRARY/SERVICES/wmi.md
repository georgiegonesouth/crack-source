# WMI

> [→ Enumeration Techniques](CHEATSHEETS/ENUMERATION/SERVICES/wmi.md)

## Overview

WMI (Windows Management Instrumentation) is Microsoft's implementation of the WBEM standard. It provides a unified interface for querying and managing virtually every aspect of a Windows system — hardware, software, processes, services, users, and the registry. WMI can be used remotely via DCOM/RPC, making it a powerful lateral movement and persistence mechanism.

## Typical Targets

- All Windows systems (WMI is built into every version of Windows)
- Domain controllers and member servers
- Workstations in enterprise environments

## Common Ports

| Port | Protocol | Usage |
|------|----------|-------|
| 135 | TCP | DCOM RPC endpoint mapper |
| Dynamic | TCP | High ports negotiated via RPC (typically 49152–65535) |

## Common Attack Paths

- **Remote Command Execution:** With valid credentials, WMI can spawn processes on a remote host via `Win32_Process.Create` — a built-in, often less-monitored alternative to PsExec.
- **Pass-the-Hash:** WMI over DCOM supports NTLM authentication, allowing hash-based authentication with Impacket's `wmiexec.py` without knowing the plaintext password.
- **WMI Event Subscriptions for Persistence:** Attackers create permanent WMI event subscriptions that survive reboots and execute payloads when a trigger condition is met.
- **Lateral Movement:** `wmiexec.py`, CrackMapExec's `--exec-method wmiexec`, and PowerShell's `Invoke-WmiMethod` are standard tools for moving laterally across Windows networks.
- **Information Gathering:** WMI queries expose running processes, installed software, logged-on users, network configuration, patch levels, and security product names.
