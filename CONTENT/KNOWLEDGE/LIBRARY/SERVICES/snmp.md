# SNMP

> [→ Enumeration Techniques](CHEATSHEETS/ENUMERATION/SERVICES/snmp.md)

## Overview

SNMP (Simple Network Management Protocol) is used to monitor and manage network devices and servers. Agents running on managed devices expose a Management Information Base (MIB) — a hierarchical tree of OIDs (Object Identifiers) containing system information. SNMP v1 and v2c use community strings as a simple shared secret; v3 adds proper authentication and encryption.

## Typical Targets

- Routers, switches, and firewalls
- Printers and UPS devices
- Servers and workstations (Windows SNMP service, net-snmp on Linux)
- Industrial and IoT devices

## Common Ports

| Port | Protocol | Usage |
|------|----------|-------|
| 161 | UDP | SNMP queries (GET, GETNEXT, SET) |
| 162 | UDP | SNMP traps (device-to-manager alerts) |

## Common Attack Paths

- **Default Community Strings:** The default read community string `public` and write string `private` are left unchanged on a large number of devices, granting full read (and sometimes write) access to the MIB.
- **Community String Brute Force:** Community strings can be brute-forced with wordlists since there is no lockout mechanism in v1/v2c.
- **Information Disclosure:** SNMP MIBs can expose running processes, installed software, network interfaces, ARP tables, routing tables, open TCP/UDP ports, and sometimes usernames.
- **SNMP SET Abuse:** With the write community string, an attacker can modify device configuration — changing interface parameters, rerouting traffic, or disabling a device.
- **SNMPv3 Brute Force:** Even with authentication, SNMPv3 user credentials can be brute-forced if the target is reachable and usernames can be enumerated.
- **SNMP Trap Interception:** Traps are sent in clear text (v1/v2c) and may contain sensitive information about device state changes, authentication failures, and errors.
