# **SMB**

> [→ Service Library](CHEATSHEETS/LIBRARY/SERVICES/smb.md)

## **SMBCLIENT**

> ### List Shares (Null Session) 
>
>     smbclient -N -L //<TARGET_IP> 

> ### Connect to a specific Share (Null Session)
>
>     smbclient -N //<TARGET_IP>/<SHARENAME>

> ### Connect to a specific Share as a User
>
>     smbclient --user=<USERNAME> //<TARGET_IP>/<SHARENAME>

> ### Flags
>
>     -N: omit password
>
>     -L: list shares
>
>     --user=username: username
>
>     --password=password: password
>
>     -I: IP-Address
>
>     -p: port

## **SMBCLIENT Console Commands**

> ### List possible Commands
>
>     help

> ### List Contents
>
>     ls

> ### Download Files 
>
>     get

> ### Upload local Files
>
>     put

> ### Execute local System Commands
>
>     !<cmd>




## RPCCLIENT

> ### Interact with the Target
>
>     rpcclient -U "" <TARGET_IP>

> ### Interact with the Target as a User
>
>     rpcclient -U <username> --password=<password> <TARGET_IP>

> ### Flags
>
>     -N: omit password
>
>     -I: IP-Address
>
>     -p: port

## RPCCLIENT Console Commands

> ### Get Server Information
>
>     srvinfo

> ### Enumerate all Network Domains
>
>     enumdomains

> ### Domain, Server and User Information
>
>     querydominfo

> ### Enumerate all available Shares
>
>     netshareenumall

> ### Provide Info on a Share
>
>     netsharegetinfo <sharename>

> ### Enumerate Domain Users
>
>     enumdomusers

> ### Provide i Info on a User
>
>     queryuser <RID>




## SMBMAP

> ### Enumerating SMB Shares
>
>     smbmap -H <TARGET_IP>

> ### Enumerating SMB Shares as a User
>
>     smbmap -u <USERNAME> -p <PASSWORD> -H <TARGET_IP>

> ### Enumerating SMB Shares on specific Port
>
>     smbmap -H <TARGET_IP> -P <PORT>

## Other Enumeration Scripts

> ### Impacket
>
>     impacket-samrdump <TARGET_IP>

> ### CrackMapExec
>
>     crackmapexec smb <TARGET_IP> --shares -u '' -p '' 

> ### enum4linux
>
>     python3 enum4linux-ng.py <TARGET_IP> -A
>
> ```
> - git clone https://github.com/cddmp/enum4linux-ng.git 
> - cd enum4linux-ng
> - pip3 install -r requirements.txt
> ```






## Tips

- Port: 445 or 137,138,139
