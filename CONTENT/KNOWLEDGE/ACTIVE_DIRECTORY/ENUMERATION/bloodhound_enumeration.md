# BloodHound Enumeration

## Collection

### SharpHound (Windows)
    .\SharpHound.exe -c All -d <DOMAIN> --zipfilename loot.zip

### BloodHound.py (Linux)
    bloodhound-python -d <DOMAIN> -u <USER> -p '<PASSWORD>' -ns <TARGET_IP> -c all

### NetExec BloodHound
    nxc ldap <TARGET_IP> -u <USER> -p '<PASSWORD>' -d <DOMAIN> --bloodhound -ns <NAMESERVER> --collection All
