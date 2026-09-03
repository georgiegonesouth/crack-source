# RPC Enumeration

## Commands

### RPCClient Null Session
    rpcclient -U "" -N <TARGET_IP>

### RPCClient Enumerate Domain Users
    rpcclient> enumdomusers

### RPCClient Enumerate Domain Groups
    rpcclient> enumdomgroups

### RPCClient Query User
    rpcclient> queryuser 0x1f4

### enum4linux-ng Full Scan
    enum4linux-ng -A <TARGET_IP>
