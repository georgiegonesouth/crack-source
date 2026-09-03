# **FTP**

> [→ Service Library](CHEATSHEETS/LIBRARY/SERVICES/ftp.md)



## **Interact with the FTP service on the target**

> ### Using FTP
>
>     ftp <TARGET_IP>

> ### Using NetCat
>
>     nc -nv <TARGET_IP> 21

> ### Using telnet
>
>     telnet <TARGET_IP> 21

> ### Using an encrypted connection
>
>     openssl s_client -connect <TARGET_IP>:21 - starttls ftpv    

    
## Login with Password 

    ftp <USER>:<PASSWORD>@<TARGET_IP>


## Download all available files on the target FTP server

    wget -m --no-passive ftp://<USER>:<PASSWORD>@<TARGET_IP> 
    
## **FTP console commands**

> ### List 
>
>     ls

> ### List recursively
>
>     ls -R

> ### Debug Mode
>
>     debug

> ### Packet trace Mode
>
>     trace

> ### Download File
>
>     get

> ### Upload local File
>
>     put



## Tips

Port: 20 or 990 for FTPS