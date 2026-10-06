# RPC

## RPCClient

### Null Session
    rpcclient -U "" -N <TARGET_IP>

### Enumerate Domain Users
    rpcclient> enumdomusers

### Enumerate Domain Groups
    rpcclient> enumdomgroups

### Query User
    rpcclient> queryuser 0x1f4

## Enum4Linux

### enum4linux-ng Full Scan
    enum4linux-ng -A <TARGET_IP>

## Change Windows User Password From Linux

> Note: 
> - this requires ForcePassword rights on the target user

### Net Rpc

    net rpc password "<TARGET_USER>" "newP@ssword2022" -U "<DOMAIN>"/"<USER>"%"<PASSWORD>" -S "<DOMAIN>"