# MEDUSA

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

## FTP

### Fixed User

    medusa -M ftp -h <TARGET_IP> -u <USER> -P <WORDLIST>

### User List

    medusa -M ftp -h <TARGET_IP> -U <USERLIST> -P <WORDLIST>

## HTTP

### Fixed User

    medusa -M http -h <TARGET_IP> -u <USER> -P <WORDLIST> -m DIR:/login.php -m FORM:username=^USER^&password=^PASS^

### User List

    medusa -M http -h <TARGET_IP> -U <USERLIST> -P <WORDLIST> -m DIR:/login.php -m FORM:username=^USER^&password=^PASS^

## IMAP

### Fixed User

    medusa -M imap -h <TARGET_IP> -u <USER> -P <WORDLIST>

### User List

    medusa -M imap -h <TARGET_IP> -U <USERLIST> -P <WORDLIST>

## MYSQL

### Fixed User

    medusa -M mysql -h <TARGET_IP> -u <USER> -P <WORDLIST>

### User List

    medusa -M mysql -h <TARGET_IP> -U <USERLIST> -P <WORDLIST>

## POP3

### Fixed User

    medusa -M pop3 -h <TARGET_IP> -u <USER> -P <WORDLIST>

### User List

    medusa -M pop3 -h <TARGET_IP> -U <USERLIST> -P <WORDLIST>

## RDP

### Fixed User

    medusa -M rdp -h <TARGET_IP> -u <USER> -P <WORDLIST>

### User List

    medusa -M rdp -h <TARGET_IP> -U <USERLIST> -P <WORDLIST>

## SSH

### Fixed User

    medusa -M ssh -h <TARGET_IP> -u <USER> -P <WORDLIST>

### User List

    medusa -M ssh -h <TARGET_IP> -U <USERLIST> -P <WORDLIST>

## SVN

### Fixed User

    medusa -M svn -h <TARGET_IP> -u <USER> -P <WORDLIST>

### User List

    medusa -M svn -h <TARGET_IP> -U <USERLIST> -P <WORDLIST>

## TELNET

### Fixed User

    medusa -M telnet -h <TARGET_IP> -u <USER> -P <WORDLIST>

### User List

    medusa -M telnet -h <TARGET_IP> -U <USERLIST> -P <WORDLIST>

## VNC

### Fixed User

    medusa -M vnc -h <TARGET_IP> -u <USER> -P <WORDLIST>

### User List

    medusa -M vnc -h <TARGET_IP> -U <USERLIST> -P <WORDLIST>

## WEB FORM

### Fixed User

    medusa -M web-form -h <TARGET_IP> -u <USER> -P <WORDLIST> -m FORM:"username=^USER^&password=^PASS^:F=Invalid"

### User List

    medusa -M web-form -h <TARGET_IP> -U <USERLIST> -P <WORDLIST> -m FORM:"username=^USER^&password=^PASS^:F=Invalid"
