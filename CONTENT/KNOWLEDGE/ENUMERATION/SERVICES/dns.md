# DNS

> [→ Service Library](CHEATSHEETS/LIBRARY/SERVICES/dns.md)

## NS request to the specific nameserver (return the nameservers of the domain)
    
    dig ns <DOMAIN> @<NAMESERVER>

## ANY request to the specific nameserver 
    
    dig any <DOMAIN> @<NAMESERVER>

## AXFR request to the specific nameserver (Zone Transer)
    
    dig axfr <DOMAIN> @<NAMESERVER>

## Subdomain brute forcing

> ### Using a specific Nameserver
>
>     dnsenum --dnsserver <NAMESERVER> --enum -p 0 -s 0 -o found_subdomains.txt -f <WORDLIST> <DOMAIN> -r

> ### Using only a Domain
>
>     dnsenum --enum <DOMAIN> -f <WORDLIST> -r

> ### good wordlists:
> <!-- token-section:WORDLIST -->
>
>```
> - /opt/useful/seclists/Discovery/DNS/namelist.txt
>```
>```
> - /opt/useful/seclists/Discovery/DNS/subdomains-top1million-110000.txt
>```
>```
> - /usr/share/seclists/Discovery/DNS/subdomains-top1million-5000.txt
>```
>```
> - /usr/share/seclists/Discovery/DNS/subdomains-top1million-20000.txt 
>```
## Tips

- Port: 53