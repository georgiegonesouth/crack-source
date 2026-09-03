# SMTP

> [→ Enumeration Techniques](CHEATSHEETS/ENUMERATION/SERVICES/smtp.md)

## Overview

SMTP (Simple Mail Transfer Protocol) is the standard protocol for sending and relaying email between mail servers and from clients to servers. It is present on virtually every mail infrastructure. While sending mail is its primary purpose, SMTP is also commonly abused for user enumeration and open relay exploitation.

## Typical Targets

- Mail transfer agents (Postfix, Exim, Sendmail, Microsoft Exchange)
- Any server hosting or relaying email
- Internet-facing mail gateways

## Common Ports

| Port | Protocol | Usage |
|------|----------|-------|
| 25 | TCP | MTA-to-MTA relay (server-to-server) |
| 587 | TCP | Submission (client-to-server, STARTTLS) |
| 465 | TCP | SMTPS (implicit TLS, legacy) |

## Common Attack Paths

- **User Enumeration (VRFY / EXPN / RCPT TO):** Many SMTP servers respond differently to valid vs. invalid addresses via the `VRFY` and `EXPN` commands, or by accepting/rejecting `RCPT TO` — enabling account enumeration without authentication.
- **Open Relay Abuse:** An SMTP server that relays mail for arbitrary external senders is an open relay, abused for spam campaigns and phishing at scale.
- **Credential Brute Force:** SMTP AUTH (SASL) can be brute-forced on submission ports (587/465) to compromise mail accounts.
- **Clear-text Credential Sniffing:** SMTP without STARTTLS transmits credentials and message content in plain text.
- **Phishing from Trusted Infrastructure:** Compromising an internal SMTP server enables sending phishing emails originating from a trusted domain, bypassing SPF/DKIM checks.
- **Mail Server CVEs:** Exim, Postfix, and Sendmail have had critical CVEs allowing RCE (e.g., Exim CVE-2019-10149 "Return of the WIZard").
