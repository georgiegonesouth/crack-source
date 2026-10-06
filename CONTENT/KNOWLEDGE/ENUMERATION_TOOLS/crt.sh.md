# **crt.sh**

## Webiste

https://crt.sh/

## Query Subdomain Names with curl
```
curl -s "https://crt.sh/?q=<DOMAIN>&output=json" | jq -r '.[]  
 | select(.name_value | contains("<SUBDOMAIN>")) | .name_value' | sort -u
```