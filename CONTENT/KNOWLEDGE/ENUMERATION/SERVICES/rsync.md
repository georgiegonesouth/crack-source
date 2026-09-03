# **RSYNC**

> [→ Service Library](CHEATSHEETS/LIBRARY/SERVICES/rsync.md)

## Probing 

    nc -nv <TARGET_IP> 873

## Enumerating Open Shares

    rsync -av --list-only rsync://<TARGET_IP>/<SHARENAME>


## Enumerating Open Shares over SSH

    rsync -e "ssh -p<PORT>" -av --list-only rsync://<TARGET_IP>/<SHARENAME>


## Tips

- Port: 873