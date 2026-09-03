# DNS

> [→ Enumeration Techniques](CHEATSHEETS/ENUMERATION/SERVICES/dns.md)

## Overview

DNS (Domain Name System) is the protocol responsible for resolving human-readable hostnames to IP addresses and vice versa. It is a fundamental part of almost every network and internet-facing infrastructure.

## Typical Targets

- Internal DNS servers (Active Directory domain controllers often double as DNS servers)
- Internet-facing authoritative nameservers
- Routers and firewalls with DNS forwarding enabled
- Any host running a recursive resolver (BIND, Unbound, dnsmasq, Windows DNS)

## Common Ports

| Port | Protocol | Usage |
|------|----------|-------|
| 53 | UDP | Standard DNS queries |
| 53 | TCP | Zone transfers, large responses (DNSSEC, >512 byte answers) |

## Common Attack Paths

- **Zone Transfer (AXFR):** Misconfigured servers allow unauthenticated zone transfers, exposing all DNS records for a domain (subdomains, internal IPs, mail servers).
- **Subdomain Enumeration:** Brute-forcing or using wordlists to discover subdomains that may expose internal or forgotten assets.
- **DNS Cache Poisoning:** Injecting forged DNS responses to redirect traffic to attacker-controlled hosts.
- **DNS Tunneling:** Encoding data in DNS queries/responses to exfiltrate data or establish C2 over DNS, bypassing firewall rules that allow DNS out.
- **Dangling DNS Records:** Subdomains pointing at deprovisioned cloud resources (e.g. old S3 buckets, Azure endpoints) that can be claimed by an attacker.
- **Internal DNS Reconnaissance:** On a compromised internal host, querying internal DNS for hostnames, service records (SRV), and mail exchangers (MX) to map the network.
