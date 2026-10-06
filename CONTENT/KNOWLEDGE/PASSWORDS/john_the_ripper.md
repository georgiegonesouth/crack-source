# JOHN THE RIPPER

## John

> ### Single Mode
>
>     john --single <PASSWORD_FILE>

> ### Wordlist Mode
>
>     john --wordlist=<WORDLIST> <PASSWORD_FILE>

> ### Incremental Mode
>
>     john --incremental <PASSWORD_FILE>

> ### Hash Format
>
>     john --format=<FORMAT> <PASSWORD_FILE>

## John Hash Formats

| Hash format | Example command | Description |
|---|---|---|
| afs | `john --format=afs` | AFS (Andrew File System) password hashes |
| bfegg | `john --format=bfegg` | bfegg hashes used in Eggdrop IRC bots |
| bf | `john --format=bf` | Blowfish-based crypt(3) hashes |
| bsdi | `john --format=bsdi` | BSDi crypt(3) hashes |
| crypt(3) | `john --format=crypt` | Traditional Unix crypt(3) hashes |
| des | `john --format=des` | Traditional DES-based crypt(3) hashes |
| dmd5 | `john --format=dmd5` | DMD5 (Dragonfly BSD MD5) password hashes |
| dominosec | `john --format=dominosec` | IBM Lotus Domino 6/7 password hashes |
| EPiServer SID hashes | `john --format=episerver` | EPiServer SID (Security Identifier) password hashes |
| hdaa | `john --format=hdaa` | hdaa password hashes used in Openwall GNU/Linux |
| hmac-md5 | `john --format=hmac-md5` | hmac-md5 password hashes |
| hmailserver | `john --format=hmailserver` | hmailserver password hashes |
| ipb2 | `john --format=ipb2` | Invision Power Board 2 password hashes |
| krb4 | `john --format=krb4` | Kerberos 4 password hashes |
| krb5 | `john --format=krb5` | Kerberos 5 password hashes |
| LM | `john --format=LM` | LM (Lan Manager) password hashes |
| lotus5 | `john --format=lotus5` | Lotus Notes/Domino 5 password hashes |
| mscash | `john --format=mscash` | MS Cache password hashes |
| mscash2 | `john --format=mscash2` | MS Cache v2 password hashes |
| mschapv2 | `john --format=mschapv2` | MS CHAP v2 password hashes |
| mskrb5 | `john --format=mskrb5` | MS Kerberos 5 password hashes |
| mssql05 | `john --format=mssql05` | MS SQL 2005 password hashes |
| mssql | `john --format=mssql` | MS SQL password hashes |
| mysql-fast | `john --format=mysql-fast` | MySQL fast password hashes |
| mysql | `john --format=mysql` | MySQL password hashes |
| mysql-sha1 | `john --format=mysql-sha1` | MySQL SHA1 password hashes |
| NETLM | `john --format=netlm` | NETLM (NT LAN Manager) password hashes |
| NETLMv2 | `john --format=netlmv2` | NETLMv2 (NT LAN Manager version 2) password hashes |
| NETNTLM | `john --format=netntlm` | NETNTLM (NT LAN Manager) password hashes |
| NETNTLMv2 | `john --format=netntlmv2` | NETNTLMv2 (NT LAN Manager version 2) password hashes |
| NEThalfLM | `john --format=nethalflm` | NEThalfLM (NT LAN Manager) password hashes |
| md5ns | `john --format=md5ns` | md5ns (MD5 namespace) password hashes |
| nsldap | `john --format=nsldap` | nsldap (OpenLDAP SHA) password hashes |
| ssha | `john --format=ssha` | ssha (Salted SHA) password hashes |
| NT | `john --format=nt` | NT (Windows NT) password hashes |
| openssha | `john --format=openssha` | OPENSSH private key password hashes |
| oracle11 | `john --format=oracle11` | Oracle 11 password hashes |
| oracle | `john --format=oracle` | Oracle password hashes |
| pdf | `john --format=pdf` | PDF (Portable Document Format) password hashes |
| phpass-md5 | `john --format=phpass-md5` | PHPass-MD5 (Portable PHP password hashing framework) password hashes |
| phps | `john --format=phps` | PHPS password hashes |
| pix-md5 | `john --format=pix-md5` | Cisco PIX MD5 password hashes |
| po | `john --format=po` | Po (Sybase SQL Anywhere) password hashes |
| rar | `john --format=rar` | RAR (WinRAR) password hashes |
| raw-md4 | `john --format=raw-md4` | Raw MD4 password hashes |
| raw-md5 | `john --format=raw-md5` | Raw MD5 password hashes |
| raw-md5-unicode | `john --format=raw-md5-unicode` | Raw MD5 Unicode password hashes |
| raw-sha1 | `john --format=raw-sha1` | Raw SHA1 password hashes |
| raw-sha224 | `john --format=raw-sha224` | Raw SHA224 password hashes |
| raw-sha256 | `john --format=raw-sha256` | Raw SHA256 password hashes |
| raw-sha384 | `john --format=raw-sha384` | Raw SHA384 password hashes |
| raw-sha512 | `john --format=raw-sha512` | Raw SHA512 password hashes |
| salted-sha | `john --format=salted-sha` | Salted SHA password hashes |
| sapb | `john --format=sapb` | SAP CODVN B (BCODE) password hashes |
| sapg | `john --format=sapg` | SAP CODVN G (PASSCODE) password hashes |
| sha1-gen | `john --format=sha1-gen` | Generic SHA1 password hashes |
| skey | `john --format=skey` | S/Key (One-time password) hashes |
| ssh | `john --format=ssh` | SSH (Secure Shell) password hashes |
| sybasease | `john --format=sybasease` | Sybase ASE password hashes |
| xsha | `john --format=xsha` | xsha (Extended SHA) password hashes |
| zip | `john --format=zip` | ZIP (WinZip) password hashes |

## Convert Protected Files to JtR Hash

> ### Find appropriate 2john tool
>
>     locate *2john* | grep <KEYWORD>

## Cracking protected Files

> ### Find appropriate 2john tool
>
>     locate *2john* | grep <KEYWORD>

*Note: run filetype2john filename.example > outfile.example*

> ### Cracking the Hash
>
>     john --wordlist=<WORDLIST> <PASSWORD_FILE>