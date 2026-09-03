# LDAP Enumeration

## Queries

### Anonymous bind
    ldapsearch -x -H ldap://<TARGET_IP> -b "DC=domain,DC=htb"

### Authenticated
    ldapsearch -x -H ldap://<TARGET_IP> -D "<USER>@<DOMAIN>" -w '<PASSWORD>' -b "DC=domain,DC=htb"

### Find users
    ldapsearch -x -H ldap://<TARGET_IP> -D "<USER>@<DOMAIN>" -w '<PASSWORD>' -b "DC=domain,DC=htb" "(objectClass=user)" sAMAccountName

### Find computers
    ldapsearch -x -H ldap://<TARGET_IP> -D "<USER>@<DOMAIN>" -w '<PASSWORD>' -b "DC=domain,DC=htb" "(objectClass=computer)" name

### NetExec LDAP Users
    nxc ldap <TARGET_IP> -u <USER> -p <PASSWORD> -d <DOMAIN> --users

### NetExec LDAP Groups
    nxc ldap <TARGET_IP> -u <USER> -p <PASSWORD> -d <DOMAIN> --groups
