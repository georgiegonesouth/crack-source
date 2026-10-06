# MySQL

> [→ Enumeration Techniques](CHEATSHEETS/ENUMERATION/SERVICES/mysql.md)

## Overview

MySQL is an open-source relational database management system, one of the most widely deployed databases in the world. It is the "M" in the LAMP/LEMP stack and backs the majority of web applications. It commonly runs on Linux but is also found on Windows.

## Typical Targets

- Web servers hosting PHP/Python/Ruby applications
- Linux and Windows application servers
- Shared hosting environments

## Common Ports

| Port | Protocol | Usage |
|------|----------|-------|
| 3306 | TCP | Default MySQL/MariaDB listener |

## Common Attack Paths

- **SQL Injection:** The most common attack path — malicious SQL injected through a web application can dump the entire database, bypass authentication, or escalate to OS-level access.
- **Credential Brute Force:** The root account or application accounts are often weakly protected; MySQL has no built-in lockout.
- **User-Defined Functions (UDF) for RCE:** By loading a malicious shared library as a UDF, an attacker with FILE privilege and access to the plugin directory can execute OS commands.
- **LOAD_FILE() / INTO OUTFILE:** The `LOAD_FILE()` function can read arbitrary files the MySQL user has OS-level read access to; `SELECT INTO OUTFILE` can write files, including webshells.
- **Exposed to Internet:** MySQL is frequently left accessible on port 3306 to the internet with weak or default credentials, making it a target for automated scanners.
- **Hash Extraction for Cracking:** The `mysql.user` table stores password hashes; dumping it allows offline cracking of all database accounts.
