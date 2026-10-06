# NFS

> [→ Service Library](CHEATSHEETS/LIBRARY/SERVICES/nfs.md)

## Show available NFS shares

    showmount -e <TARGET_IP>

## Mount the specific NFS share to ./target-NFS

    mount -t nfs <TARGET_IP>:/<share> ./target-NFS/ -o nolock 

## Unmount the specific NFS share

    umount ./target-NFS

## Nmap for NFS

    sudo nmap --script nfs* <TARGET_IP> -sV -p111,2049

## Tips

- Port: 111 and 2049