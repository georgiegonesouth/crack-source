# SNMP

> [→ Service Library](CHEATSHEETS/LIBRARY/SERVICES/snmp.md)

## Querying OIDs using snmpwalk

    snmpwalk -v2c -c <COMMUNITY_STRING> <TARGET_IP>

<!-- token-options:WORDLIST
/opt/useful/seclists/Discovery/SNMP/snmp.txt
-->

## Bruteforcing community strings of the SNMP service

    onesixtyone -c <WORDLIST> <TARGET_IP>
       
> ### good wordlists:
> <!-- token-section:WORDLIST -->
>
>```
> - /opt/useful/seclists/Discovery/SNMP/snmp.txt
>```


## Bruteforcing SNMP service OIDs

    braa <COMMUNITY_STRING>@<TARGET_IP>:.1.*

## Tips

- Port: 161/udp and 162/udp