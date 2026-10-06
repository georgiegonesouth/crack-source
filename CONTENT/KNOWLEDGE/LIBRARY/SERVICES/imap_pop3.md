# IMAP & POP3

> [→ Enumeration Techniques](CHEATSHEETS/ENUMERATION/SERVICES/imap_pop3.md)

## Overview

IMAP (Internet Message Access Protocol) and POP3 (Post Office Protocol v3) are used by mail clients to retrieve email from a mail server. IMAP allows managing mail on the server (folders, flags), while POP3 downloads and typically deletes mail. Both can run over plain TCP or wrapped in TLS.

## Typical Targets

- Corporate mail servers (Exchange, Postfix, Dovecot, Zimbra)
- Webmail backends
- Any host running an MDA (Mail Delivery Agent)

## Common Ports

| Port | Protocol | Usage |
|------|----------|-------|
| 143 | TCP | IMAP (plain) |
| 993 | TCP | IMAPS (TLS) |
| 110 | TCP | POP3 (plain) |
| 995 | TCP | POP3S (TLS) |

## Common Attack Paths

- **Credential Brute Force:** Mail accounts are high-value targets. Many servers allow unlimited login attempts, making brute force and password spraying viable.
- **Clear-text Credential Sniffing:** On unencrypted IMAP/POP3 (ports 143/110), credentials and mail content are transmitted in plain text.
- **Email Harvesting:** Once authenticated, an attacker can read all emails in the mailbox — useful for gathering credentials, internal documents, and social engineering material.
- **Phishing Pivot:** Access to a mail account enables sending phishing emails from a trusted internal address, bypassing many email security controls.
- **User Enumeration:** Some implementations respond differently to valid vs. invalid usernames during authentication, enabling account enumeration.
