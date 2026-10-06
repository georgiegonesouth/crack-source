# Oracle TNS

> [→ Enumeration Techniques](CHEATSHEETS/ENUMERATION/SERVICES/oracle_tns.md)

## Overview

Oracle TNS (Transparent Network Substrate) is the proprietary networking protocol used by Oracle Database for client-server communication. The TNS Listener acts as a gatekeeper that routes incoming connections to the appropriate Oracle database instance. Oracle is common in large enterprise and financial environments.

## Typical Targets

- Large enterprises and financial institutions
- Government and healthcare organizations
- Any environment running Oracle EBS, Oracle Forms, or Oracle-backed ERP systems

## Common Ports

| Port | Protocol | Usage |
|------|----------|-------|
| 1521 | TCP | Default Oracle TNS Listener |
| 1522–1529 | TCP | Alternative listener ports |
| 2483 / 2484 | TCP | Oracle TNS with native encryption |

## Common Attack Paths

- **SID Brute Force:** The Oracle Service Identifier (SID) must be known to connect. Default and common SIDs (`ORCL`, `XE`, `DB`) can be brute-forced; revealing a valid SID is the first step.
- **Credential Brute Force:** Default Oracle accounts (`system`, `sys`, `scott`, `dbsnmp`) with default passwords (`change_on_install`, `tiger`, `manager`) are frequently left in place.
- **TNS Poison Attack:** On older versions, the TNS Listener can be redirected with `ALTER SYSTEM` to proxy connections through an attacker, enabling MitM credential capture.
- **SYSDBA Privilege Abuse:** Connecting as SYSDBA (the highest Oracle privilege level) bypasses most access controls and enables full database control and OS interaction via `DBMS_SCHEDULER` or `UTL_FILE`.
- **OS Command Execution:** Through packages like `DBMS_SCHEDULER`, `DBMS_JAVA`, or external procedures, a privileged Oracle user can execute OS commands on the database host.
