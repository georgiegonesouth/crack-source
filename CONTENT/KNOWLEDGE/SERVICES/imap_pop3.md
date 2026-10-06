# IMAP & POP3

> [→ Service Library](CHEATSHEETS/LIBRARY/SERVICES/imap_pop3.md)

## Log in to the IMAPS service using cURL

    curl -k 'imaps://<TARGET_IP>' --user <user>:<password>

## Connect to the IMAP/S service

> ### Encrypted
>
>     openssl s_client -connect <TARGET_IP>:imaps

> ### Unencrypted
>
>     telnet <TARGET_IP> 143

## Connect to the POP3/S service

> ### Encrypted
>
>     openssl s_client -connect <TARGET_IP>:pop3s

> ### Unencrypted
>
>     telnet <TARGET_IP> 110

## IMAP Commands

| Command | Description |
|---------|-------------|
| `1 LOGIN username password` | User's login. |
| `1 LIST "" *` | Lists all directories. |
| `1 CREATE "INBOX"` | Creates a mailbox with a specified name. |
| `1 DELETE "INBOX"` | Deletes a mailbox. |
| `1 RENAME "ToRead" "Important"` | Renames a mailbox. |
| `1 LSUB "" *` | Returns a subset of names from the set of names that the User has declared as being active or subscribed. |
| `1 SELECT INBOX` | Selects a mailbox so that messages in the mailbox can be accessed. |
| `1 UNSELECT INBOX` | Exits the selected mailbox. |
| `1 FETCH <ID> all` | Retrieves data associated with a message in the mailbox. |
| `1 FETCH <ID> (BODY[])` | Retrieves the body of the message. |
| `1 CLOSE` | Removes all messages with the Deleted flag set. |
| `1 LOGOUT` | Closes the connection with the IMAP server. |

## POP3 Commands

| Command | Description |
|---------|-------------|
| `USER username` | Identifies the user. |
| `PASS password` | Authentication of the user using its password. |
| `STAT` | Requests the number of saved emails from the server. |
| `LIST` | Requests from the server the number and size of all emails. |
| `RETR id` | Requests the server to deliver the requested email by ID. |
| `DELE id` | Requests the server to delete the requested email by ID. |
| `CAPA` | Requests the server to display the server capabilities. |
| `RSET` | Requests the server to reset the transmitted information. |
| `QUIT` | Closes the connection with the POP3 server. |



## Tips

- Port Imap/s: 143/993
- Port Pop3/s: 110/995
