# **RDP**

> [→ Service Library](CHEATSHEETS/LIBRARY/SERVICES/rdp.md)

## Logging into RDP

    xfreerdp /u:<USER> /p:"<PASSWORD>" /v:<TARGET_IP>

## RDP Security Check
```    
    sudo cpan
```
```
    git clone https://github.com/CiscoCXSecurity/rdp-sec-check.git && cd rdp-sec-check
```
```
    ./rdp-sec-check.pl <TARGET_IP>
```

## Nmap for RDP

    nmap -sV -sC <TARGET_IP> -p3389 --script rdp*

## Tips

- Port: 3389