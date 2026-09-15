# ENCRYPTED FILES & ARCHIVES

## Find Files with Extensions
```
for ext in $(echo ".xls .xls* .xltx .od* .doc .doc* .pdf .pot .pot* .pp*");do echo -e "\nFile extension: " $ext; find / -name *$ext 2>/dev/null | grep -v "lib\|fonts\|share\|core" ;done
```

## Find Private Keys
```
grep -rnE '^\-{5}BEGIN [A-Z0-9]+ PRIVATE KEY\-{5}$' /* 2>/dev/null
```

## OpenSSL Encrypted GZIP Files
```
for i in $(cat rockyou.txt);do openssl enc -aes-256-cbc -d -in GZIP.gzip -k $i 2>/dev/null| tar xz;done
```

## Cracking BitLocker-encrypted Drives

> ### Create Hashfile
>
>     bitlocker2john -i <FILE>.vhd > backup.hashes; grep "bitlocker\$0" backup.hashes > backup.hash

> ### Crack with Hashcat 
>
>     hashcat -a 0 -m 22100 backup.hash /usr/share/wordlists/rockyou.txt

