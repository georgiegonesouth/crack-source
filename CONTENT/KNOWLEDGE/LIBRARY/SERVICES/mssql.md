# MSSQL

> [→ Enumeration Techniques](CHEATSHEETS/ENUMERATION/SERVICES/mssql.md)

## Overview

MSSQL (Microsoft SQL Server) is Microsoft's relational database management system, widely used in Windows enterprise environments. It integrates tightly with Active Directory and Windows authentication, and includes powerful built-in features that can be abused for lateral movement and command execution.

## Typical Targets

- Windows servers in enterprise environments
- Internal application backends (.NET apps, SharePoint, Dynamics)
- Domain-joined servers where the SQL service runs as a privileged AD account

## Common Ports

| Port | Protocol | Usage |
|------|----------|-------|
| 1433 | TCP | Default MSSQL listener |
| 1434 | UDP | SQL Server Browser service (instance discovery) |

## Common Attack Paths

- **xp_cmdshell for RCE:** The `xp_cmdshell` stored procedure executes OS commands as the SQL Server service account. It is disabled by default but can be re-enabled by a sysadmin.
- **SA Account Brute Force:** The built-in `sa` (system administrator) account using SQL authentication is often left enabled with a weak password.
- **Linked Server Exploitation:** SQL Servers can be configured to trust and query other SQL Servers. Chaining linked server queries can allow pivoting across database servers, including cross-domain.
- **UNC Path Injection:** Triggering `xp_dirtree` or `xp_fileexist` with a UNC path to an attacker-controlled SMB server captures the MSSQL service account's NTLMv2 hash for offline cracking or relay.
- **Credential Theft:** Service accounts running MSSQL often have elevated AD privileges; compromising the DB may directly yield domain privilege escalation.
- **Impersonation Chains:** `EXECUTE AS LOGIN` can be chained to impersonate other logins and escalate from a low-privilege SQL user to sysadmin.
