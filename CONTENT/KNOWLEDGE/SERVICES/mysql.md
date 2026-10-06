# **MySQL**

> [→ Service Library](CHEATSHEETS/LIBRARY/SERVICES/mysql.md)

## Connect to the MySQL Server

    mysql -u <USER> -p<PASSWORD> -h <TARGET_IP>

> ### IF: *ERROR 2026 (HY000): TLS/SSL error: self-signed certificate in certificate chain*
>
>     mysql -u <USER> -p<PASSWORD> -h <TARGET_IP> --skip-ssl

## MySQL Shell commands

| COMMAND | USAGE |
|---------|-------|
| `show databases;` | Show all databases. |
| `use <database>;` | Select one of the existing databases. |
| `show tables;` | Show all available tables in the selected database. |
| `show columns from <table>;` | Show all columns in the selected table. |
| `select * from <table>;` | Show everything in the desired table. |
| `select * from <table> where <column> = "<string>";` | Search for needed string in the desired table. |

## Nmap for MySQL

    sudo nmap <TARGET_IP> -sV -sC -p3306 --script mysql*

## Tips

- Port: 3306