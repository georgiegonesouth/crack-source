# LDAP

## NXC

### Find Users

    nxc ldap <TARGET_IP> -u <USER> -p <PASSWORD> -d <DOMAIN> --users

### Find Groups

    nxc ldap <TARGET_IP> -u <USER> -p <PASSWORD> -d <DOMAIN> --groups

### Find computers

    ldapsearch -x -H ldap://<TARGET_IP> -D "<USER>@<DOMAIN>" -w '<PASSWORD>' -b "DC=domain,DC=htb" "(objectClass=computer)" name

### Machine Account Quota (Can this account add other Accounts?)

    nxc ldap <IP> -u <USER> -p <PASSWORD> -M maq

## LdapSearch

### Anonymous bind
    ldapsearch -x -H ldap://<TARGET_IP> -b "DC=domain,DC=htb"

### Authenticated
    ldapsearch -x -H ldap://<TARGET_IP> -D "<USER>@<DOMAIN>" -w '<PASSWORD>' -b "DC=domain,DC=htb"

### Find users
    ldapsearch -x -H ldap://<TARGET_IP> -D "<USER>@<DOMAIN>" -w '<PASSWORD>' -b "DC=domain,DC=htb" "(objectClass=user)" sAMAccountName



