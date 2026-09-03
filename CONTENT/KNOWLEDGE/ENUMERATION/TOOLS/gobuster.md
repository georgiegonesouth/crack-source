# **GOBUSTER**

## Directory Brute Force

    gobuster dir -u http://<DOMAIN> -w <WORDLIST>

### good wordlists

> <!-- token-section:WORDLIST -->
>
>```
> - /usr/share/seclists/Discovery/Web-Content/common.txt
>```
>```
> - /usr/share/seclists/Discovery/Web-Content/raft-medium-directories.txt
>```

## Subdomain Brute Force

    gobuster dns -d <DOMAIN> -t 50 -w <WORDLIST>

### good wordlists

> <!-- token-section:WORDLIST -->
>
>```
> - /usr/share/seclists/Discovery/DNS/subdomains-top1million-5000.txt
>```
>```
> - /usr/share/seclists/Discovery/DNS/subdomains-top1million-20000.txt
>```
>```
> - /opt/useful/seclists/Discovery/DNS/subdomains-top1million-110000.txt
>```

## Vhost Brute Force

    gobuster vhost -u http://<DOMAIN> -w <WORDLIST> --append-domain

### good wordlists

> <!-- token-section:WORDLIST -->
>
>```
> - /usr/share/seclists/Discovery/DNS/subdomains-top1million-5000.txt
>```
>```
> - /usr/share/seclists/Discovery/DNS/subdomains-top1million-20000.txt
>```
>```
> - /opt/useful/seclists/Discovery/DNS/subdomains-top1million-110000.txt
>```