# **SSH**

> [→ Service Library](CHEATSHEETS/LIBRARY/SERVICES/ssh.md)

## Logging into SSH

> ### With a Password
>
>     ssh <USER>@<TARGET_IP>

> ### With a Key
>
>     ssh -i key.file <USER>@<TARGET_IP>

> ### Change Authentication Method
>
>```
>     ssh -v <USER>@<TARGET_IP> 
>```
>         - look for "Authentications that can continue: publickey,password,keyboard-interactive"
>```
>     ssh -v <USER>@<TARGET_IP> -o PreferredAuthentication=<option>
>```

## SSH Audit

This tool audits the configuration of an SSH server or client and highlights the areas needing improvement.

> ### Download
>
>     git clone https://github.com/jtesta/ssh-audit.git && cd ssh-audit

> ### Usage
>
>     ./ssh-audit.py <TARGET_IP>

## Generate an SSH Keypair

> ### id_ed25519
>
>     ssh-keygen -t ed25519 -a 100 -f ~/.ssh/id_ed25519

> ### id_rsa
>
>     ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa



## Tips
 
- Port: 22