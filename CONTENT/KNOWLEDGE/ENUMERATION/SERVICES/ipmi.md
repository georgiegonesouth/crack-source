# IPMI

> [→ Service Library](CHEATSHEETS/LIBRARY/SERVICES/ipmi.md)

## Metasploit Modules

> ### IPMI version detection
>
>     scanner/ipmi/ipmi_version

> ### Dump IPMI hashes
>
>     scanner/ipmi/ipmi_dumphashes

## Default Credentials across Vendors

| PRODUCT | USER | PASSWORD |
|---------|------|----------|
| Dell iDRAC | `root` | `calvin` |
| HP iLO | `Administrator` | randomized 8-character string of numbers and uppercase letters |
| Supermicro IPMI | `ADMIN` | `ADMIN` |

## Hashcat command for HP iLO PW Hash

```hashcat -m 7300 ipmi.txt -a 3 ?1?1?1?1?1?1?1?1 -1 ?d?u```



## Tips
- Port: 623/udp