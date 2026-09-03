## Questions

What can we see?
What reasons can we have for seeing it?
What image does what we see create for us?
What do we gain from it?
How can we use it?
What can we not see?
What reasons can there be that we do not see?
What image results for us from what we do not see?

## Principles

1.	There is more than meets the eye. Consider all points of view.

2.	Distinguish between what we see and what we do not see.

3.	There are always ways to gain more information. Understand the target.

## COPY PASTE COMMANDS

### Show Subdomains of a Target from crt.sh

    curl -s https://crt.sh/\?q\=<TARGET_IP>\&output\=json | jq . | grep name | cut -d":" -f2 | grep -v "CN=" | cut -d'"' -f2 | awk '{gsub(/\\n/,"\n");}1;' | sort -u

### Show only Company Hosted Servers

    for i in $(cat subdomainlist);do host $i | grep "has address" | grep <TARGET_IP> | cut -d" " -f1,4;done

### Shodan IP List

    for i in $(cat subdomainlist);do host $i | grep "has address" | grep <TARGET_IP> | cut -d" " -f4 >> ip-addresses.txt;done
    for i in $(cat ip-addresses.txt);do shodan host $i;done

### Display all available DNS records

    dig any <example.com>

