# SMB

> [→ Enumeration Techniques](CHEATSHEETS/ENUMERATION/SERVICES/smb.md)

## Overview

SMB (Server Message Block) is the primary Windows protocol for file sharing, printer sharing, and inter-process communication. It is the backbone of Windows networking and Active Directory. Samba implements SMB on Linux/Unix. SMB has been the source of some of the most impactful vulnerabilities in Windows history.

## Typical Targets

- Windows servers and workstations (all versions)
- Domain controllers (SYSVOL, NETLOGON shares)
- Linux/Unix systems running Samba
- NAS devices with Windows file-sharing enabled

## Common Ports

| Port | Protocol | Usage |
|------|----------|-------|
| 445 | TCP | SMB over TCP (modern, direct) |
| 139 | TCP | SMB over NetBIOS |
| 137–138 | UDP | NetBIOS Name Service / Datagram |

## Common Attack Paths

- **EternalBlue (MS17-010):** Pre-auth RCE in SMBv1, weaponized by WannaCry and NotPetya. Still found on unpatched Windows 7 / Server 2008 systems.
- **Null Session Enumeration:** Older SMB configurations allow unauthenticated (null) sessions to enumerate shares, users, groups, and password policies.
- **Pass-the-Hash (PtH):** NTLM authentication allows using a password hash directly without knowing the plaintext. Tools like `smbclient`, CrackMapExec, and Impacket support PtH.
- **NTLM Relay Attacks:** Capturing NTLM authentication attempts (via Responder) and relaying them to other hosts allows authenticating as the victim. Enables lateral movement without cracking hashes.
- **Coercion Attacks (PrintSpooler / PetitPotam / DFSCoerce):** These force a Windows host to authenticate to an attacker-controlled server, capturing its machine account hash for relay or cracking.
- **Share Enumeration:** Readable shares frequently contain scripts, config files, and backup files with plaintext credentials, SSH keys, or sensitive documents.
- **GPP Credential Disclosure:** Old Group Policy Preferences XML files in SYSVOL may contain AES-encrypted (but easily decrypted) passwords (`cpassword` attribute).
