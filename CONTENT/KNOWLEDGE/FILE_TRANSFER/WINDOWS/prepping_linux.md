# **PREPPING LINUX**

> [→ Windows Commands](CHEATSHEETS/FILE_TRANSFER/WINDOWS/windows_file_transfer.md)

## Create SMB Server

> ### Create Unauthenticated Server on Linux
>
>     sudo impacket-smbserver share -smb2support <DIRECTORY>

> ### Create Authenticated Server on Linux
>
>```
> net use n: \\<IP>\share /user:test test
>```
>
>```
> copy n:\<FILE>
>```


## Create FTP Server

> ### Install FTP Server module 
>
>     sudo pip3 install pyftpdlib

> ### Setting up a Python3 FTP Server
>
>     sudo python3 -m pyftpdlib --port 21

> ### Setting up a Python3 FTP Server, allow Uploads
>
>     sudo python3 -m pyftpdlib --port 21 --write