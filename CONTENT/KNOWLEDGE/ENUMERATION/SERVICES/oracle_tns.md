# Oracle TNS

> [→ Service Library](CHEATSHEETS/LIBRARY/SERVICES/oracle_tns.md)

## ODAT.PY

> ### Perform a Scan
>
>     ./odat.py all -s <TARGET_IP>

### Upload a file with Oracle RDBMS
>
>     ./odat.py utlfile -s <TARGET_IP> -d <DATABASE> -U <USER> -P <PASSWORD> --sysdba --putFile C:\\insert\\path file.txt ./file.txt

> ### Download Odat.py
>```
> sudo apt-get update
> sudo apt-get install -y build-essential python3-dev libaio1
> cd ~
> wget https://files.pythonhosted.org/packages/source/c/cx_Oracle/cx_Oracle-8.3.0.tar.gz
> tar xzf cx_Oracle-8.3.0.tar.gz
> cd cx_Oracle-8.3.0
> python3 setup.py build
> sudo python3 setup.py install
> cd ~
> git clone https://github.com/quentinhardy/odat.git
> cd odat/
> pip install python-libnmap
> git submodule init
> git submodule update
> sudo apt-get install python3-scapy -y
> sudo pip3 install colorlog termcolor passlib python-libnmap
> sudo apt-get install build-essential libgmp-dev -y
> pip3 install pycryptodome
> pip3 install openpyxl
>```

## SQLPlus

> ### Log in to the Oracle database / as admin
>
>```
>     sqlplus <USER>/<PASSWORD>@<TARGET_IP>/<DATABASE> 
>```
>```
>     sqlplus <USER>/<PASSWORD>@<TARGET_IP>/<DATABASE> as sysdba
>```

> ### Download SQLPlus
>
>```
> sudo apt update
> sudo apt upgrade parrot-core
> sudo apt update
> sudo apt install oracle-instantclient-sqlplus
>```
>

> ### IF *sqlplus: error while loading shared libraries: libsqlplus.so: cannot open shared object file: No such file or directory*
>
>```
> sudo sh -c "echo /usr/lib/oracle/12.2/client64/lib > /etc/ld.so.conf.d/oracle-instantclient.conf";sudo ldconfig
>```

## ORACLE TNS COMMANDS

| COMMAND | USAGE |
|---------|-------|
| `select table_name from all_tables;` | Show all tables. |
| `select * from user_role_privs;` | Show current user privileges. |
| `select name, password from sys.user$;` | Show users and password hashes. |



## Tips

- Port: 1521