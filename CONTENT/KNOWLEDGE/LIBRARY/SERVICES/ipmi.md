# IPMI

> [→ Enumeration Techniques](CHEATSHEETS/ENUMERATION/SERVICES/ipmi.md)

## Overview

IPMI (Intelligent Platform Management Interface) is a standardized hardware-level management interface built into server motherboards. It allows out-of-band management (power control, console access, sensor monitoring) independently of the OS, even when the server is powered off. It is implemented via a Baseboard Management Controller (BMC) — sold under vendor names like iDRAC (Dell), iLO (HPE), and IPMI (Supermicro).

## Typical Targets

- Bare-metal servers in data centers
- Enterprise server infrastructure (Dell, HPE, Supermicro, Lenovo)
- Servers on isolated management networks (OOB/OOBM VLANs)

## Common Ports

| Port | Protocol | Usage |
|------|----------|-------|
| 623 | UDP | IPMI / RMCP+ (main IPMI port) |
| 80 / 443 | TCP | BMC web interface |
| 22 | TCP | BMC SSH interface |

## Common Attack Paths

- **RAKP Authentication Hash Disclosure (CVE-2013-4786):** In IPMI 2.0, the server sends a salted HMAC hash of the password during the authentication handshake even before the user authenticates. This hash can be captured and cracked offline.
- **Cipher 0 Authentication Bypass:** When Cipher Suite 0 is enabled, IPMI accepts any password for any user — effectively allowing unauthenticated access.
- **Default Credentials:** BMC interfaces commonly ship with well-known default credentials (e.g., `ADMIN:ADMIN`, `root:calvin`, `Administrator:changeme`) that are rarely changed.
- **BMC Web Interface Vulnerabilities:** Vendor-specific vulnerabilities in the web management interface (e.g., HPE iLO 4 RCE, Supermicro buffer overflows) can provide root-level OS access.
- **Privilege Escalation via BMC OS Access:** Compromising the BMC gives full control over the host server — including power cycling, KVM console, and direct memory access.
