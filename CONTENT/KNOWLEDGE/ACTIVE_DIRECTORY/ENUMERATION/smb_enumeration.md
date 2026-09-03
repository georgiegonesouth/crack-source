# SMB Enumeration

## Commands

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
