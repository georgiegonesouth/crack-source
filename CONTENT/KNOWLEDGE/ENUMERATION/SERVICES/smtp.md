# SMTP

> [→ Service Library](CHEATSHEETS/LIBRARY/SERVICES/smtp.md)

## Interact with the target 

    telnet <TARGET_IP> 25

## Commands

| COMMAND | USAGE |
|---------|-------|
| `AUTH PLAIN` | AUTH is a service extension used to authenticate the client. |
| `HELO` | The client logs in with its computer name and thus starts the session. |
| `MAIL FROM` | The client names the email sender. |
| `RCPT TO` | The client names the email recipient. |
| `DATA` | The client initiates the transmission of the email. |
| `RSET` | The client aborts the initiated transmission but keeps the connection between client and server. |
| `VRFY` | The client checks if a mailbox is available for message transfer. |
| `EXPN` | The client also checks if a mailbox is available for messaging with this command. |
| `NOOP` | The client requests a response from the server to prevent disconnection due to time-out. |
| `QUIT` | The client terminates the session. |

<!-- token-options:WORDLIST
/usr/share/wordlists/metasploit/unix_users.txt
-->

## User enumeration script

    smtp-user-enum -M VRFY -U <WORDLIST> -t <TARGET_IP> -m 60 -w 20
> ### good wordlists:
> <!-- token-section:WORDLIST -->
>
>```
> - /usr/share/wordlists/metasploit/unix_users.txt
>```
        
## Check if server is an open relay

    sudo nmap <TARGET_IP> -p25 --script smtp-open-relay -v



## Tips

- Port: 25